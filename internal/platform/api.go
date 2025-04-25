package platform

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum"
	ethereumClient "github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/client"
	"github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/pollers"
	ethereumRepository "github.com/jeronimobarea/transaction_parser/internal/chains/ethereum/repository"
	"github.com/jeronimobarea/transaction_parser/internal/parser"
	parserHandlers "github.com/jeronimobarea/transaction_parser/internal/parser/handlers"
	"github.com/jeronimobarea/transaction_parser/pkg/httpx"
	"github.com/jeronimobarea/transaction_parser/pkg/osx"
)

func Run(ctx context.Context) {
	logger := log.Default()

	var ethereumParser parser.Parser
	{

		ethereumNodePRCUrl := osx.GetEnvFallback("ETHEREUM_NODE_RPC_URL", "https://ethereum-rpc.publicnode.com")
		ethClient := ethereumClient.NewClient(ethereumNodePRCUrl, logger)
		ethereumRepo := ethereumRepository.NewMemoryStorage()

		ethereumParser = ethereum.NewEthereumParser(ethereumRepo, logger)

		poller := pollers.NewPoller(ethClient, ethereumRepo, logger)

		runner := pollers.NewRunner(logger, 5*time.Second)

		go func() {
			err := runner.Run(ctx, poller)
			if err != nil {
				logger.Fatal(err)
			}
		}()
	}

	var parserSvc parser.Service
	{
		parserSvc = parser.NewService(logger)

		//** Register parsers **//
		parserSvc.Register(parser.EthereumChainID, ethereumParser)
	}

	var router *httpx.Router
	{
		router = httpx.NewRouter()
	}
	parserHandlers.RegisterRoutes(router, parserSvc, logger)

	httpServerPort := osx.GetEnvFallback("HTTP_SERVER_PORT", ":3000")
	logger.Printf("Server listening on: %s \n", httpServerPort)
	http.ListenAndServe(httpServerPort, router)
}
