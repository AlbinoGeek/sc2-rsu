MAKE                = make --no-print-directory

DESCRIBE           := $(shell git describe --match "v*" --always --tags)
DESCRIBE_PARTS     := $(subst -, ,$(DESCRIBE))

VERSION_TAG        := $(word 1,$(DESCRIBE_PARTS))
COMMITS_SINCE_TAG  := $(word 2,$(DESCRIBE_PARTS))

VERSION            := $(subst v,,$(VERSION_TAG))
VERSION_PARTS      := $(subst ., ,$(VERSION))

MAJOR              := $(word 1,$(VERSION_PARTS))
MINOR              := $(word 2,$(VERSION_PARTS))
MICRO              := $(word 3,$(VERSION_PARTS))

NEXT_MAJOR         := $(shell echo $$(($(MAJOR)+1)))
NEXT_MINOR         := $(shell echo $$(($(MINOR)+1)))
NEXT_MICRO          = $(shell echo $$(($(MICRO)+$(COMMITS_SINCE_TAG))))

BINARYNAME          = sc2-rsu
MODNAME             = github.com/AlbinoGeek/sc2-rsu
TARGETDIR           = _dist

ifeq ($(strip $(COMMITS_SINCE_TAG)),)
CURRENT_VERSION_MICRO := $(MAJOR).$(MINOR).$(MICRO)
CURRENT_VERSION_MINOR := $(CURRENT_VERSION_MICRO)
CURRENT_VERSION_MAJOR := $(CURRENT_VERSION_MICRO)
else
CURRENT_VERSION_MICRO := $(MAJOR).$(MINOR).$(NEXT_MICRO)
CURRENT_VERSION_MINOR := $(MAJOR).$(NEXT_MINOR).0
CURRENT_VERSION_MAJOR := $(NEXT_MAJOR).0.0
endif

DATE                = $(shell date +'%d.%m.%Y')
TIME                = $(shell date +'%H:%M:%S')
COMMIT             := $(shell git rev-parse HEAD)
AUTHOR             := $(firstword $(subst @, ,$(shell git show --format="%aE" $(COMMIT))))
BRANCH_NAME        := $(shell git rev-parse --abbrev-ref HEAD)

TAG_MESSAGE         = "$(TIME) $(DATE) $(AUTHOR) $(BRANCH_NAME)"
COMMIT_MESSAGE     := $(shell git log --format=%B -n 1 $(COMMIT))

CURRENT_TAG_MICRO  := "v$(CURRENT_VERSION_MICRO)"
CURRENT_TAG_MINOR  := "v$(CURRENT_VERSION_MINOR)"
CURRENT_TAG_MAJOR  := "v$(CURRENT_VERSION_MAJOR)"

# --- Version commands ---

.PHONY: version
version:
	@$(MAKE) version-micro

.PHONY: version-micro
version-micro:
	@echo "$(CURRENT_VERSION_MICRO)"

.PHONY: version-minor
version-minor:
	@echo "$(CURRENT_VERSION_MINOR)"

.PHONY: version-major
version-major:
	@echo "$(CURRENT_VERSION_MAJOR)"

# --- Tag commands ---

.PHONY: tag-micro
tag-micro:
	@echo "$(CURRENT_TAG_MICRO)"

.PHONY: tag-minor
tag-minor:
	@echo "$(CURRENT_TAG_MINOR)"

.PHONY: tag-major
tag-major:
	@echo "$(CURRENT_TAG_MAJOR)"

# --- Meta info ---

.PHONY: tag-message
tag-message:
	@echo "$(TAG_MESSAGE)"

.PHONY: commit-message
commit-message:
	@echo "$(COMMIT_MESSAGE)"

# --- Actual Make Targets ---

.PHONY: all
all: "$(TARGETDIR)/$(BINARYNAME)"

"$(TARGETDIR)/$(BINARYNAME)":
	if [ ! -d "$(TARGETDIR)" ]; then mkdir "$(TARGETDIR)"; fi
	go build -ldflags "-X main.PROGRAM=$(BINARYNAME) -X main.VERSION=$(CURRENT_VERSION_MICRO)" -o "$(TARGETDIR)/$(BINARYNAME)" "$(MODNAME)" # -s -w

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.PROGRAM=$(BINARYNAME) -X main.VERSION=$(CURRENT_VERSION_MICRO)" -o "$(TARGETDIR)/$(BINARYNAME)-$(CURRENT_VERSION_MICRO)-linux-amd64" "$(MODNAME)"
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.PROGRAM=$(BINARYNAME) -X main.VERSION=$(CURRENT_VERSION_MICRO)" -o "$(TARGETDIR)/$(BINARYNAME)-$(CURRENT_VERSION_MICRO)-windows-amd64.exe" "$(MODNAME)"
	GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.PROGRAM=$(BINARYNAME) -X main.VERSION=$(CURRENT_VERSION_MICRO)" -o "$(TARGETDIR)/$(BINARYNAME)-$(CURRENT_VERSION_MICRO)-linux-386" "$(MODNAME)"
	GOOS=windows GOARCH=386 go build -ldflags "-s -w -X main.PROGRAM=$(BINARYNAME) -X main.VERSION=$(CURRENT_VERSION_MICRO)" -o "$(TARGETDIR)/$(BINARYNAME)-$(CURRENT_VERSION_MICRO)-windows-386.exe" "$(MODNAME)"

.PHONY: clean
clean:
	rm -f "$(TARGETDIR)/$(BINARYNAME)"

.PHONY: update
update:
	go get -u -v ./...
