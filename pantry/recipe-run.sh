#!/bin/bash

# Recipe Crawler Run Script
# This script helps you run the recipe crawler with optimal settings

# Default settings
WORKERS=20
DEPTH=4
DELAY=1.5
MAX_REQUESTS=10
DEBUG=true

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Print banner
echo -e "${GREEN}=======================================${NC}"
echo -e "${GREEN}      Recipe Crawler Run Script       ${NC}"
echo -e "${GREEN}=======================================${NC}"

# Function to show help
show_help() {
    echo -e "\n${YELLOW}Usage:${NC}"
    echo -e "  ./run.sh [command] [options]"
    echo -e "\n${YELLOW}Commands:${NC}"
    echo -e "  crawl [url]     Crawl recipes (optionally from a specific URL)"
    echo -e "  test [url]      Test a specific URL for recipe data"
    echo -e "  delete          Delete the Elasticsearch index"
    echo -e "  help            Show this help message"
    echo -e "\n${YELLOW}Options:${NC}"
    echo -e "  -w, --workers N       Number of concurrent workers (default: $WORKERS)"
    echo -e "  -d, --depth N         Maximum crawl depth (default: $DEPTH)"
    echo -e "  -t, --delay N         Delay between requests in seconds (default: $DELAY)"
    echo -e "  -m, --max-requests N  Maximum concurrent requests per domain (default: $MAX_REQUESTS)"
    echo -e "  --debug               Enable debug mode (default: $DEBUG)"
    echo -e "  --no-debug            Disable debug mode"
    echo -e "\n${YELLOW}Examples:${NC}"
    echo -e "  ./run.sh crawl"
    echo -e "  ./run.sh crawl https://www.allrecipes.com/recipes/"
    echo -e "  ./run.sh crawl -w 30 -d 5 -t 2"
    echo -e "  ./run.sh test https://www.allrecipes.com/recipe/8000/award-winning-chocolate-chip-cookies/"
    echo -e "\n"
}

# Parse command
if [ $# -eq 0 ]; then
    show_help
    exit 0
fi

COMMAND=$1
shift

# Special case for help
if [ "$COMMAND" == "help" ]; then
    show_help
    exit 0
fi

# Parse options
URL=""
while [[ $# -gt 0 ]]; do
    case $1 in
        -w|--workers)
            WORKERS="$2"
            shift 2
            ;;
        -d|--depth)
            DEPTH="$2"
            shift 2
            ;;
        -t|--delay)
            DELAY="$2"
            shift 2
            ;;
        -m|--max-requests)
            MAX_REQUESTS="$2"
            shift 2
            ;;
        --debug)
            DEBUG=true
            shift
            ;;
        --no-debug)
            DEBUG=false
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            # If it's not an option, assume it's a URL
            if [[ $1 == http* ]]; then
                URL="$1"
            fi
            shift
            ;;
    esac
done

# Execute command
case $COMMAND in
    crawl)
        echo -e "${GREEN}Starting recipe crawler...${NC}"
        echo -e "${YELLOW}Workers:${NC} $WORKERS"
        echo -e "${YELLOW}Max Depth:${NC} $DEPTH"
        echo -e "${YELLOW}Delay:${NC} $DELAY seconds"
        echo -e "${YELLOW}Max Requests Per Domain:${NC} $MAX_REQUESTS"
        echo -e "${YELLOW}Debug Mode:${NC} $DEBUG"

        CMD="go run *.go index -workers=$WORKERS -depth=$DEPTH -delay=$DELAY -max-requests=$MAX_REQUESTS -debug=$DEBUG"

        # If URL is provided, add it
        if [ ! -z "$URL" ]; then
            echo -e "${YELLOW}Starting URL:${NC} $URL"
            CMD="$CMD $URL"
        else
            echo -e "${YELLOW}Starting URLs:${NC} Default recipe sites"
        fi

        echo -e "\n${GREEN}Running Command:${NC} $CMD"
        echo -e "${GREEN}=======================================${NC}"

        # Execute the command
        $CMD
        ;;

    test)
        if [ -z "$URL" ]; then
            echo -e "${RED}Error: No URL provided for testing${NC}"
            echo -e "Usage: ./run.sh test [url]"
            exit 1
        fi

        echo -e "${GREEN}Testing URL: $URL${NC}"
        go run *.go test-url $URL
        ;;

    delete)
        echo -e "${RED}Warning: This will delete all recipe data!${NC}"
        read -p "Are you sure you want to continue? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${GREEN}Deleting Elasticsearch index...${NC}"
            go run *.go delete
        else
            echo -e "${YELLOW}Operation cancelled${NC}"
        fi
        ;;

    *)
        echo -e "${RED}Error: Unknown command '$COMMAND'${NC}"
        show_help
        exit 1
        ;;
esac

echo -e "\n${GREEN}Done!${NC}"