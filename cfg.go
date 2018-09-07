package botbase

type Config struct {
    TGBot struct {
        Token string
    }

    Proxy_SOCKS5 struct {
        Server string
        User string
        Pass string
    }

    Redis struct {
        Server string
        DB int
        Pass string
    }
}
