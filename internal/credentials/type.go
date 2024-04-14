package credentials

type Credentials struct {
	Alias      string  `msgpack:"alias"`
	Username   string  `msgpack:"userame"`
	Host       string  `msgpack:"host"`
	Port       int     `msgpack:"port"`
	Password   *string `msgpack:"password,omitempty"`
	Key        []byte  `msgpack:"key,omitempty"`
	Passphrase []byte  `msgpack:"passphrase,omitempty"`
	Source     string  `msgpack:"-"`
}
