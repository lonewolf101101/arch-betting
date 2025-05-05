# Variables
CC = go run
MODE = -mode=debug
FRONTEND_DIR = frontend
BACKEND_DIR = backend

# Determine the processor type
PROCESSOR_TYPE := $(shell uname -m)

web:
ifeq ($(PROCESSOR_TYPE), arm64)
	cd $(BACKEND_DIR)/cmd && $(CC) -tags=dynamic ./web $(MODE)
else
	cd $(BACKEND_DIR)/cmd && $(CC) ./web $(MODE)
endif

ui: 
	cd ./frontend && yarn dev --dotenv ./env/.env --host

superset-saver:
	cd $(BACKEND_DIR)/cmd && $(CC) -tags=dynamic ./superset-saver $(MODE)