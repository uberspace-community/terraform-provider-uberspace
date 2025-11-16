package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/internal/provider"
)

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	apiKey := os.Getenv("UBERSPACE_APIKEY")
	if apiKey == "" {
		return fmt.Errorf("UBERSPACE_APIKEY environment variable is not set")
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("please provide at least one asteroid name as argument")
	}

	for _, asteroid := range os.Args[1:] {
		fmt.Println("Resetting asteroid:", asteroid)

		c, err := client.NewClient("https://marvin.uberspace.is", client.WithClient(provider.NewAuthClient(apiKey)))
		if err != nil {
			return fmt.Errorf("failed to create Uberspace client: %w", err)
		}

		domains, err := c.AsteroidsWebdomainsList(ctx, client.AsteroidsWebdomainsListParams{
			AsteroidName: asteroid,
		})
		if err != nil {
			return fmt.Errorf("failed to list webdomains: %w", err)
		}

		for _, domain := range domains.Results {
			backends, err := c.AsteroidsWebdomainsBackendsList(ctx, client.AsteroidsWebdomainsBackendsListParams{
				AsteroidName:  domain.Asteroid,
				WebdomainName: domain.Domain,
			})
			if err != nil {
				return fmt.Errorf("failed to list backends for domain %s: %w", domain.Domain, err)
			}

			for _, backend := range backends.Results {
				fmt.Printf("Deleting backend %s for domain %s\n", backend.Path, domain.Domain)
				if err := c.AsteroidsWebdomainsBackendsDelete(ctx, client.AsteroidsWebdomainsBackendsDeleteParams{
					Path:          backend.Path,
					AsteroidName:  backend.Asteroid,
					WebdomainName: domain.Domain,
				}); err != nil {
					return fmt.Errorf("failed to delete backend %s for domain %s: %w", backend.Path, domain.Domain, err)
				}
				fmt.Printf("Deleted backend %s for domain %s\n", backend.Path, domain.Domain)
			}

			headers, err := c.AsteroidsWebdomainsHeadersList(ctx, client.AsteroidsWebdomainsHeadersListParams{
				AsteroidName:  domain.Asteroid,
				WebdomainName: domain.Domain,
			})
			if err != nil {
				return fmt.Errorf("failed to list headers for domain %s: %w", domain.Domain, err)
			}

			for _, header := range headers.Results {
				fmt.Printf("Deleting header %s for domain %s\n", header.Name, domain.Domain)
				if err := c.AsteroidsWebdomainsHeadersDelete(ctx, client.AsteroidsWebdomainsHeadersDeleteParams{
					ID:            strconv.Itoa(header.Pk),
					AsteroidName:  domain.Asteroid,
					WebdomainName: domain.Domain,
				}); err != nil {
					return fmt.Errorf("failed to delete header %s for domain %s: %w", header.Name, domain.Domain, err)
				}
				fmt.Printf("Deleted header %s for domain %s\n", header.Name, domain.Domain)
			}

			if domain.Domain == fmt.Sprintf("%s.uber8.space", asteroid) {
				fmt.Printf("Skipping deletion of primary webdomain %s\n", domain.Domain)
				continue
			}

			fmt.Printf("Deleting webdomain %s\n", domain.Domain)
			if err := c.AsteroidsWebdomainsDelete(ctx, client.AsteroidsWebdomainsDeleteParams{
				AsteroidName: domain.Asteroid,
				Name:         domain.Domain,
			}); err != nil {
				return fmt.Errorf("failed to delete webdomain %s: %w", domain.Domain, err)
			}
			fmt.Printf("Deleted webdomain %s\n", domain.Domain)
		}

		mailDomains, err := c.AsteroidsMaildomainsList(ctx, client.AsteroidsMaildomainsListParams{
			AsteroidName: asteroid,
		})
		if err != nil {
			return fmt.Errorf("failed to list maildomains: %w", err)
		}

		for _, domain := range mailDomains.Results {
			backends, err := c.AsteroidsMaildomainsUsersList(ctx, client.AsteroidsMaildomainsUsersListParams{
				AsteroidName:   domain.Asteroid,
				MaildomainName: domain.Domain,
			})
			if err != nil {
				return fmt.Errorf("failed to list backends for maildomain %s: %w", domain.Domain, err)
			}

			for _, user := range backends.Results {
				fmt.Printf("Deleting user %s for maildomain %s\n", user.Name, domain.Domain)
				if err := c.AsteroidsMaildomainsUsersDelete(ctx, client.AsteroidsMaildomainsUsersDeleteParams{
					Local:          user.Name,
					AsteroidName:   domain.Asteroid,
					MaildomainName: domain.Domain,
				}); err != nil {
					return fmt.Errorf("failed to delete user %s for maildomain %s: %w", user.Name, domain.Domain, err)
				}
				fmt.Printf("Deleted user %s for maildomain %s\n", user.Name, domain.Domain)
			}

			if domain.Domain == fmt.Sprintf("%s.uber8.space", asteroid) {
				fmt.Printf("Skipping deletion of primary maildomain %s\n", domain.Domain)
				continue
			}

			fmt.Printf("Deleting maildomain %s\n", domain.Domain)
			if err := c.AsteroidsMaildomainsDelete(ctx, client.AsteroidsMaildomainsDeleteParams{
				AsteroidName: domain.Asteroid,
				Name:         domain.Domain,
			}); err != nil {
				return fmt.Errorf("failed to delete maildomain %s: %w", domain.Domain, err)
			}
			fmt.Printf("Deleted maildomain %s\n", domain.Domain)
		}
	}

	return nil
}
