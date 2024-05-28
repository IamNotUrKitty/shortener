package workers

import (
	"context"
	"fmt"

	"github.com/iamnoturkkitty/shortener/internal/app/links/handlers"
	"github.com/iamnoturkkitty/shortener/internal/domain/links"
)

func DeleteLinksWorker(queue <-chan links.DeleteLinkTask, repo handlers.Repository) {
	var buffer []links.DeleteLinkTask
	for {
		select {
		case linkTask := <-queue:
			buffer = append(buffer, linkTask)
			if len(buffer) > 10 {
				err := repo.DeleteLinkBatch(context.Background(), buffer)
				if err != nil {
					fmt.Println("panic")
				}

				buffer = nil
			}

		default:
			if len(buffer) > 0 {
				err := repo.DeleteLinkBatch(context.Background(), buffer)
				if err != nil {
					fmt.Println("panic")
				}

				buffer = nil
			}
		}
	}
}
