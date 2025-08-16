package credentials

type Credentials struct {
	Name       string  `msgpack:"name"`
	UserName   string  `msgpack:"username"`
	Host       string  `msgpack:"host"`
	Port       int     `msgpack:"port"`
	Password   *string `msgpack:"password,omitempty"`
	Key        []byte  `msgpack:"key,omitempty"`
	Passphrase []byte  `msgpack:"passphrase,omitempty"`
}
