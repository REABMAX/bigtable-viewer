.PHONY: init,run

init:
	BIGTABLE_EMULATOR_HOST=localhost:8086 PROJECT=project INSTANCE=bigtable-instance go run cmd/init/main.go

run:
	BIGTABLE_EMULATOR_HOST=localhost:8086 PROJECT=project INSTANCE=bigtable-instance go run cmd/web/main.go