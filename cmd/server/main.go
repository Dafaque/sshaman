package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Dafaque/sshaman/internal/server/auth"
	"github.com/Dafaque/sshaman/internal/server/config"
	"github.com/Dafaque/sshaman/internal/server/controllers/roles"
	"github.com/Dafaque/sshaman/internal/server/controllers/users"
	"github.com/Dafaque/sshaman/internal/server/db"
	_ "github.com/Dafaque/sshaman/internal/server/db/migrations"
	"github.com/Dafaque/sshaman/internal/server/handler"
	rolesRepo "github.com/Dafaque/sshaman/internal/server/repositories/roles"
	usersRepo "github.com/Dafaque/sshaman/internal/server/repositories/users"
	api "github.com/Dafaque/sshaman/pkg/server/api"
)

var configFile = flag.String("config", "config.yaml", "Path to the configuration file")

func main() {
	// MARK: - logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// MARK: - config
	flag.Parse()
	cfg, err := config.New(*configFile)
	if err != nil {
		logger.Fatal("failed to parse config", zap.Error(err))
	}

	// MARK: - db
	db, err := db.New(cfg)
	if err != nil {
		logger.Fatal("failed to create db", zap.Error(err))
	}
	defer db.Close() //@todo gracefully shutdown

	// MARK: - migrations
	goose.SetDialect("postgres")
	goose.SetLogger(newGooseLogger(logger))
	err = goose.Up(db, ".")
	if err != nil {
		logger.Fatal("failed to migrate", zap.Error(err))
	}

	// MARK: - repositories
	repos, err := newRepositories(db)
	if err != nil {
		logger.Fatal("failed to create repositories", zap.Error(err))
	}

	// MARK: - controllers
	usersController, err := users.New(repos.users)
	if err != nil {
		logger.Fatal("failed to create users controller", zap.Error(err))
	}

	rolesController, err := roles.New(repos.roles)
	if err != nil {
		logger.Fatal("failed to create roles controller", zap.Error(err))
	}
	// MARK: - auth
	jwtManager, err := auth.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.Issuer,
	)
	if err != nil {
		logger.Fatal("failed to create jwt manager", zap.Error(err))
	}
	interceptor := auth.NewGRPCAuthInterceptor(jwtManager, rolesController, usersController)

	// MARK: - post init controllers
	if usersController.PrintToken() {
		tok, err := jwtManager.GenerateToken(users.SUID)
		if err != nil {
			logger.Fatal("failed to generate token for superuser", zap.Error(err))
		}
		logger.Info("superuser token", zap.String("token", tok))
	}

	// MARK: - grpc
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	handler := handler.New(
		usersController,
		rolesController,
	)
	api.RegisterRemoteCredentialsManagerServer(s, handler)
	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close() //@todo gracefully shutdown

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type gooseLogger struct {
	logger *zap.Logger
}

func newGooseLogger(logger *zap.Logger) *gooseLogger {
	return &gooseLogger{logger: logger.Named("goose")}
}
func (l *gooseLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal("fatal", zap.String("format", fmt.Sprintf(format, v...)))
}
func (l *gooseLogger) Printf(format string, v ...interface{}) {
	l.logger.Info("info", zap.String("format", fmt.Sprintf(format, v...)))
}

type repositories struct {
	users usersRepo.Repository
	roles rolesRepo.Repository
}

func newRepositories(db *sql.DB) (*repositories, error) {
	users := usersRepo.New(db)
	roles := rolesRepo.New(db)
	return &repositories{
		users: users,
		roles: roles,
	}, nil
}