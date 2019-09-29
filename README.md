
Marshal & Unmarshal go to js.Value

## Example

```go
type addr struct {
    Street string
    Num    int
}
type user struct {
    Name    string
    Email   []string
    IsAdmin bool
    Limit   int
    Addr    *addr
}

// create new user in go
usr := &user{
    Name:    "admin",
    Email:   []string{"admin@domain", "root@domain"},
    IsAdmin: true,
    Limit:   9001,
    Addr: &addr{
        Street: "somewhere",
        Num:    99,
    },
}

// marshal to js.Value
usrjs, err := jsref.Marshal(usr)
if err != nil {
    return err
}

// get values from js and print
fmt.Println("JS Values")
fmt.Println("---------")
fmt.Println("Name", usrjs.Get("Name").String())
fmt.Println("IsAdmin", usrjs.Get("IsAdmin").Bool())
fmt.Println("Limit", usrjs.Get("Limit").Int())
for i := 0; i < usrjs.Get("Email").Length(); i++ {
    fmt.Println("Email", i, usrjs.Get("Email").Index(i))
}
fmt.Println("Addr.Street", usrjs.Get("Addr").Get("Street").String())
fmt.Println("Addr.Num", usrjs.Get("Addr").Get("Num").Int())
fmt.Println("---------")

// set values from js
usrjs.Set("Name", "user")
usrjs.Set("IsAdmin", false)
usrjs.Set("Limit", 100)
usrjs.Set("Email", []interface{}{"user@domain"})
usrjs.Get("Addr").Set("Street", "somewhere else")
usrjs.Get("Addr").Set("Num", 2)

usr = &user{}

// unmarshal js to user struct
err = jsref.Unmarshal(usr, usrjs)
if err != nil {
    return err
}

// get and print user values 
fmt.Println("Go Values")
fmt.Println("---------")
fmt.Println("Name", usr.Name)
fmt.Println("IsAdmin", usr.IsAdmin)
fmt.Println("Limit", usr.Limit)
for i, e := range usr.Email {
    fmt.Println("Email", i, e)
}
fmt.Println("Addr.Street", usr.Addr.Street)
fmt.Println("Addr.Num", usr.Addr.Num)
fmt.Println("---------")
```

Prints out

```shell
JS Values
---------
Name admin
IsAdmin true
Limit 9001
Email 0 admin@domain
Email 1 root@domain
Addr.Street somewhere
Addr.Num 99
---------
Go Values
---------
Name user
IsAdmin false
Limit 100
Email 0 user@domain
Addr.Street somewhere else
Addr.Num 2
---------
```

You can run the test by issuing:

`GOOS=js GOARCH=wasm go test -v -run TestExample`

More on how to setup your system to run tests with webassembly can be found here: [running-tests-in-the-browser](https://github.com/golang/go/wiki/WebAssembly#running-tests-in-the-browser)

## Tags

Apply `jsref:"<val>"` tag to a struct field to alter default behavior

### Ignore

```go
type MyType {
    Name string
    Password string `jsref:"-"`
}
```

### Rename

```go
type MyType {
    UserName string `jsref:"userName"`
    Password string `jsref:"password"`
}
```

## Purpose

Implemented originally for having a convenient way to attach structured data in the `details` field when creating a custom event and subsequently retrieve them at the event handler.

See: [Creating_and_triggering_events](https://developer.mozilla.org/en-US/docs/Web/Guide/Events/Creating_and_triggering_events)

## TODO

JsrefMarshaler & JsrefUnmarshaler interfaces for custom parsing