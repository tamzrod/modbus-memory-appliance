package main

func appMain() {
	cfg := loadConfig()
	memories := buildMemories(cfg)
	ingestSvc := buildIngest(memories)

	startModbus(cfg, memories)
	startMQTT(cfg, ingestSvc)
	startREST(cfg, memories, ingestSvc)
	startRawIngest(cfg, memories)


	blockForever()
}
