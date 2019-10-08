# Logging entry point lib for Butik services

## Fast start
1. Select formatter:
    ```shell script
       formatter := &TextFormatter{
         TextFormatter: logrus.TextFormatter{ForceColors: true, FullTimestamp: true},
       }
       // Or this:
       formatter := &JSONFormatter{}
    ```
2. Initialize logging instance:
    ```shell script
       goerrors.InitLog(goerrors.DebugLevel, formatter, nil)
    ```
3. Use as logrus:
    ```shell script
       goerrors.Log().Info("test")
    ```

For detailed usage see code examples and comments in `goerrors_test.go`

## Levels
TRACE, DEBUG, INFO --> stdout

WARN, ERROR, FATAL, PANIC --> stderr + Sentry

## Requests logging and recovery
The library provides middleware for GRPC and HTTP servers:
```shell script
    r := mux.Router()
    r.Use(goerrors.HTTPLoggingHandler) // INFO-level short info for every incoming request
    r.Use(goerrors.HTTPRecoverer)
``` 

```shell script
    grpcServer := grpc.NewServer(
        grpc.StreamInterceptor(GRPCStreamServerRecoverer()),
    	grpc.UnaryInterceptor(GRPCUnaryServerRecoverer()),
    )
```