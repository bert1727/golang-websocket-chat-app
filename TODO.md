# TODO

добавить это
// Загружаем историю через Axios (используя наш инстанс с интерсепторами)
api.get(`/messages/${receiverId}`).then(res => {
setMessages(res.data);
}).catch(console.error);

- Add tests for methods
- different mods for prod and dev app

- [x] Create return types for services and handlers
- [x] Add relations between messages and users via foreign keys in user table
- [x] Add validation for models
(partially)
- [x] Move all middleware in separate file
- [x] Add logout endpoint (add it on client)

```
 // consoleWriter := zerolog.ConsoleWriter{
 //  Out:        os.Stdout,
 //  NoColor:    false,
 //  TimeFormat: "24",
 // }
 // logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
 // log.Logger = zerolog.New(consoleWriter).With().Caller().Timestamp().Caller().Logger()
 // app.Use(middleware.New(middleware.Config{
 //  Logger: &logger,
 // }))

 // logger := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
 // app.Use(logger.New(logger.Config{
 //  Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
 // }))
```
