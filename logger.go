package main

import (
    "github.com/sirupsen/logrus"
    "os"
)

// Initialize logger
func InitLogger() *logrus.Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})
    file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
        logger.SetOutput(file)
    } else {
        logger.Info("Failed to log to file, using default stderr")
    }
    return logger
}
