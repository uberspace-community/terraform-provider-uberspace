//nolint:forbidigo
package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

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

	c, err := client.NewClient("https://marvin.uberspace.is", client.WithClient(provider.NewAuthClient(apiKey)))
	if err != nil {
		return fmt.Errorf("failed to create Uberspace client: %w", err)
	}

	for _, asteroid := range os.Args[1:] {
		fmt.Println("Resetting asteroid:", asteroid)

		if err := resetWebdomain(ctx, c, asteroid); err != nil {
			return err
		}

		if err := resetMaildomain(ctx, c, asteroid); err != nil {
			return err
		}
	}

	return nil
}

func resetWebdomain(ctx context.Context, c *client.Client, asteroid string) error {
	domains, err := c.AsteroidsWebdomainsList(ctx, client.AsteroidsWebdomainsListParams{
		AsteroidName: asteroid,
	})
	if err != nil {
		return fmt.Errorf("failed to list webdomains: %w", err)
	}

	for _, domain := range domains.Results {
		if err := resetWebdomainBackend(ctx, c, domain); err != nil {
			return err
		}

		if err := resetWebdomainHeader(ctx, c, domain); err != nil {
			return err
		}

		if domain.Name == fmt.Sprintf("%s.uber.space", asteroid) {
			fmt.Printf("Skipping deletion of primary webdomain %s\n", domain.Name)

			continue
		}

		fmt.Printf("Deleting webdomain %s\n", domain.Name)

		if err := c.AsteroidsWebdomainsDelete(ctx, client.AsteroidsWebdomainsDeleteParams{
			AsteroidName: domain.Asteroid,
			Name:         domain.Name,
		}); err != nil {
			return fmt.Errorf("failed to delete webdomain %s: %w", domain.Name, err)
		}

		fmt.Printf("Deleted webdomain %s\n", domain.Name)
	}

	return nil
}

func resetWebdomainBackend(ctx context.Context, c *client.Client, domain client.WebDomain) error {
	backends, err := c.AsteroidsWebdomainsBackendsList(ctx, client.AsteroidsWebdomainsBackendsListParams{
		AsteroidName:  domain.Asteroid,
		WebdomainName: domain.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to list backends for domain %s: %w", domain.Name, err)
	}

	for _, backend := range backends.Results {
		fmt.Printf("Deleting backend %s for domain %s\n", backend.Path, domain.Name)

		if err := c.AsteroidsWebdomainsBackendsDelete(ctx, client.AsteroidsWebdomainsBackendsDeleteParams{
			Path:          backend.Path,
			AsteroidName:  backend.Asteroid,
			WebdomainName: domain.Name,
		}); err != nil {
			return fmt.Errorf("failed to delete backend %s for domain %s: %w", backend.Path, domain.Name, err)
		}

		fmt.Printf("Deleted backend %s for domain %s\n", backend.Path, domain.Name)
	}

	return nil
}

func resetWebdomainHeader(ctx context.Context, c *client.Client, domain client.WebDomain) error {
	headers, err := c.AsteroidsWebdomainsHeadersList(ctx, client.AsteroidsWebdomainsHeadersListParams{
		AsteroidName:  domain.Asteroid,
		WebdomainName: domain.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to list headers for domain %s: %w", domain.Name, err)
	}

	for _, header := range headers.Results {
		fmt.Printf("Deleting header %s for domain %s\n", header.Name, domain.Name)

		if err := c.AsteroidsWebdomainsHeadersDelete(ctx, client.AsteroidsWebdomainsHeadersDeleteParams{
			ID:            strconv.Itoa(header.Pk),
			AsteroidName:  domain.Asteroid,
			WebdomainName: domain.Name,
		}); err != nil {
			return fmt.Errorf("failed to delete header %s for domain %s: %w", header.Name, domain.Name, err)
		}

		fmt.Printf("Deleted header %s for domain %s\n", header.Name, domain.Name)
	}

	return nil
}

func resetMaildomain(ctx context.Context, c *client.Client, asteroid string) error {
	mailDomains, err := c.AsteroidsMaildomainsList(ctx, client.AsteroidsMaildomainsListParams{
		AsteroidName: asteroid,
	})
	if err != nil {
		return fmt.Errorf("failed to list maildomains: %w", err)
	}

	for _, domain := range mailDomains.Results {
		backends, err := c.AsteroidsMaildomainsUsersList(ctx, client.AsteroidsMaildomainsUsersListParams{
			AsteroidName:   domain.Asteroid,
			MaildomainName: domain.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to list backends for maildomain %s: %w", domain.Name, err)
		}

		for _, user := range backends.Results {
			if slices.Contains([]string{"sysmail", "postmaster", "abuse", "hostmaster"}, user.Name) {
				fmt.Printf("Skipping deletion of system user %s for maildomain %s\n", user.Name, domain.Name)

				continue
			}

			fmt.Printf("Deleting user %s for maildomain %s\n", user.Name, domain.Name)

			if err := c.AsteroidsMaildomainsUsersDelete(ctx, client.AsteroidsMaildomainsUsersDeleteParams{
				Local:          user.Name,
				AsteroidName:   domain.Asteroid,
				MaildomainName: domain.Name,
			}); err != nil {
				return fmt.Errorf("failed to delete user %s for maildomain %s: %w", user.Name, domain.Name, err)
			}

			fmt.Printf("Deleted user %s for maildomain %s\n", user.Name, domain.Name)
		}

		systemDomains := []string{
			fmt.Sprintf("%s.uber.space", asteroid),
			fmt.Sprintf("mail.%s.uber.space", asteroid),
		}

		if slices.Contains(systemDomains, domain.Name) || (strings.HasPrefix(domain.Name, "mail-") && strings.HasSuffix(domain.Name, fmt.Sprintf(".%s.uber.space", asteroid))) {
			fmt.Printf("Skipping deletion of primary maildomain %s\n", domain.Name)

			continue
		}

		fmt.Printf("Deleting maildomain %s\n", domain.Name)

		if err := c.AsteroidsMaildomainsDelete(ctx, client.AsteroidsMaildomainsDeleteParams{
			AsteroidName: domain.Asteroid,
			Name:         domain.Name,
		}); err != nil {
			return fmt.Errorf("failed to delete maildomain %s: %w", domain.Name, err)
		}

		fmt.Printf("Deleted maildomain %s\n", domain.Name)
	}

	return nil
}
