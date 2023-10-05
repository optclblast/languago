package test

import (
	"languago/internal/test/generators"
	"languago/pkg/logger"
	"testing"
)

func TestLoggers(t *testing.T) {
	var dbg bool = true

	// std logger

	log1 := logger.NewDefaultLogger(dbg)

	log1.Debug("STD DEBUG", generators.NewPairs())
	//log.Err("STD ERR", generators.NewPairs())
	log1.Info("STD INFO", generators.NewPairs())
	log1.Log("STD LOG", generators.NewPairs())
	log1.Warn("STD WARN", generators.NewPairs())

	// zap

	log2 := logger.NewZapWrapper(dbg, logger.EnvParam_LOCAL)

	log2.Debug("ZAP DEBUG", generators.NewPairs())
	//log.Err("ZAP ERR", generators.NewPairs())
	log2.Info("ZAP INFO", generators.NewPairs())
	log2.Log("ZAP LOG", generators.NewPairs())
	log2.Warn("ZAP WARN", generators.NewPairs())

	// logrus

	log3 := logger.NewLogrusWrapper(dbg, logger.EnvParam_LOCAL)

	log3.Debug("LOGRUS DEBUG", generators.NewPairs())
	//log.Err("LOGRUS ERR", generators.NewPairs())
	log3.Info("LOGRUS INFO", generators.NewPairs())
	log3.Log("LOGRUS LOG", generators.NewPairs())
	log3.Warn("LOGRUS WARN", generators.NewPairs())

}

func TestLogrus(t *testing.T) {
	var dbg bool = true

	log3 := logger.NewLogrusWrapper(dbg, logger.EnvParam_LOCAL)

	log3.Debug("LOGRUS DEBUG", generators.NewPairs())
	//log.Err("LOGRUS ERR", generators.NewPairs())
	log3.Info("LOGRUS INFO", generators.NewPairs())
	log3.Log("LOGRUS LOG", generators.NewPairs())
	log3.Warn("LOGRUS WARN", generators.NewPairs())

}

func TestZap(t *testing.T) {
	var dbg bool = true

	log2 := logger.NewZapWrapper(dbg, logger.EnvParam_LOCAL)

	log2.Debug("ZAP DEBUG", generators.NewPairs())
	//log.Err("ZAP ERR", generators.NewPairs())
	log2.Info("ZAP INFO", generators.NewPairs())
	log2.Log("ZAP LOG", generators.NewPairs())
	log2.Warn("ZAP WARN", generators.NewPairs())
}

func TestSTD(t *testing.T) {
	var dbg bool = true

	log1 := logger.NewDefaultLogger(dbg)

	log1.Debug("STD DEBUG", generators.NewPairs())
	//log.Err("STD ERR", generators.NewPairs())
	log1.Info("STD INFO", generators.NewPairs())
	log1.Log("STD LOG", generators.NewPairs())
	log1.Warn("STD WARN", generators.NewPairs())
}
