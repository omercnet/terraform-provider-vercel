package vercel_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func getFirewallImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("no ID is set")
		}

		if rs.Primary.Attributes["team_id"] == "" {
			return rs.Primary.ID, nil
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["team_id"], rs.Primary.Attributes["project_id"]), nil
	}
}

func TestAcc_FirewallConfigResource(t *testing.T) {
	name := acctest.RandString(16)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallConfigResource(name, teamIDConfig(t)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.xss.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.sqli.action",
						"log"),
					resource.TestCheckNoResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.sqli.active"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.rce.active",
						"false"),
					resource.TestCheckNoResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.php"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"enabled",
						"true"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.name",
						"test"),
					resource.TestCheckResourceAttrWith(
						"vercel_firewall_config.custom",
						"rules.rule.0.id",
						func(rule_id string) error {
							if !strings.HasPrefix(rule_id, "rule_test") {
								return fmt.Errorf("expected id does not match got %s - expected %s", rule_id, "rule_test_...")
							}
							return nil
						}),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.action.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.condition_group.0.conditions.0.type",
						"path"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.condition_group.0.conditions.1.value",
						"POST"),
					resource.TestCheckResourceAttrWith(
						"vercel_firewall_config.custom",
						"rules.rule.1.id",
						func(rule_id string) error {
							if !strings.HasPrefix(rule_id, "rule_test2") {
								return fmt.Errorf("expected id does not match got %s - expected %s", rule_id, "rule_test2_...")
							}
							return nil
						}),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.action",
						"rate_limit"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.limit",
						"100"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.algo",
						"fixed_window"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.action",
						"challenge"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.window",
						"60"),
					resource.TestCheckTypeSetElemAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.keys.*",
						"ip"),
					resource.TestCheckTypeSetElemAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.keys.*",
						"ja4"),
					resource.TestCheckResourceAttrWith(
						"vercel_firewall_config.custom",
						"rules.rule.2.id",
						func(rule_id string) error {
							if !strings.HasPrefix(rule_id, "rule_test3") {
								return fmt.Errorf("expected id does not match got %s - expected %s", rule_id, "rule_test3_...")
							}
							return nil
						}),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.2.action.redirect.location",
						"/bye"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.2.action.redirect.permanent",
						"false"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.4.condition_group.0.conditions.0.values.0",
						"/test1"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.4.condition_group.0.conditions.0.values.1",
						"/test2"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.0.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.0.hostname",
						"test.com"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.1.ip",
						"1.2.3.4/32"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.1.hostname",
						"*.test.com"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.2.ip",
						"2.4.6.8"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.2.hostname",
						"*"),

					resource.TestCheckResourceAttr(
						"vercel_firewall_config.botfilter",
						"managed_rulesets.bot_filter.action",
						"challenge"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.botfilter",
						"managed_rulesets.bot_filter.active",
						"true"),
				),
			},
			{
				ImportState:       true,
				ResourceName:      "vercel_firewall_config.managed",
				ImportStateIdFunc: getFirewallImportID("vercel_firewall_config.managed"),
			},
			{
				ImportState:       true,
				ResourceName:      "vercel_firewall_config.custom",
				ImportStateIdFunc: getFirewallImportID("vercel_firewall_config.custom"),
			},
			{
				ImportState:       true,
				ResourceName:      "vercel_firewall_config.ips",
				ImportStateIdFunc: getFirewallImportID("vercel_firewall_config.ips"),
			},
			{
				ImportState:       true,
				ResourceName:      "vercel_firewall_config.botfilter",
				ImportStateIdFunc: getFirewallImportID("vercel_firewall_config.botfilter"),
			},
			{
				Config: testAccFirewallConfigResourceUpdated(name, teamIDConfig(t)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.xss.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.sqli.action",
						"deny"),
					resource.TestCheckNoResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.sqli.active"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.rce.active",
						"false"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.php.action",
						"log"),
					resource.TestCheckNoResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.php.active"),
					resource.TestCheckNoResourceAttr(
						"vercel_firewall_config.managed",
						"managed_rulesets.owasp.java"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"enabled",
						"true"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.name",
						"test1"),
					resource.TestCheckResourceAttrWith(
						"vercel_firewall_config.custom",
						"rules.rule.0.id",
						func(rule_id string) error {
							if !strings.HasPrefix(rule_id, "rule_test1") {
								return fmt.Errorf("expected id does not match got %s - expected %s", rule_id, "rule_test1_...")
							}
							return nil
						}),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.action.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.condition_group.0.conditions.0.type",
						"path"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.0.condition_group.0.conditions.1.value",
						"POST"),
					resource.TestCheckResourceAttrWith(
						"vercel_firewall_config.custom",
						"rules.rule.1.id",
						func(rule_id string) error {
							if !strings.HasPrefix(rule_id, "rule_test2") {
								return fmt.Errorf("expected id does not match got %s - expected %s", rule_id, "rule_test2_...")
							}
							return nil
						}),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.action",
						"rate_limit"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.limit",
						"150"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.algo",
						"fixed_window"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.action",
						"challenge"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.window",
						"60"),
					resource.TestCheckTypeSetElemAttr(
						"vercel_firewall_config.custom",
						"rules.rule.1.action.rate_limit.keys.*",
						"ip"),
					resource.TestCheckResourceAttrWith(
						"vercel_firewall_config.custom",
						"rules.rule.2.id",
						func(rule_id string) error {
							if !strings.HasPrefix(rule_id, "rule_test3") {
								return fmt.Errorf("expected id does not match got %s - expected %s", rule_id, "rule_test3_...")
							}
							return nil
						}),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.2.action.redirect.location",
						"/bye"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.custom",
						"rules.rule.2.action.redirect.permanent",
						"false"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.neg",
						"rules.rule.0.condition_group.0.conditions.0.neg",
						"true"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.neg",
						"rules.rule.0.condition_group.0.conditions.0.values.0",
						"1.2.3.4"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.neg",
						"rules.rule.0.condition_group.0.conditions.0.values.1",
						"3.4.5.6"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.0.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.0.hostname",
						"test.com"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.1.ip",
						"1.2.3.4/32"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.1.hostname",
						"*.test.com"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.2.ip",
						"2.4.6.8/32"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.ips",
						"ip_rules.rule.2.hostname",
						"*"),

					resource.TestCheckResourceAttr(
						"vercel_firewall_config.botfilter",
						"managed_rulesets.bot_filter.action",
						"deny"),
					resource.TestCheckResourceAttr(
						"vercel_firewall_config.botfilter",
						"managed_rulesets.bot_filter.active",
						"false"),
				),
			},
		},
	})
}

func testAccFirewallConfigResource(name, teamID string) string {
	return fmt.Sprintf(`
resource "vercel_project" "managed" {
    name = "test-acc-%[1]s-mrs"
    %[2]s
}

resource "vercel_firewall_config" "managed" {
    project_id = vercel_project.managed.id
    %[2]s

    managed_rulesets {
        owasp {
            xss = { action = "deny" }
            sqli = { action = "log" }

            rce = { action = "deny", active = false }
        }
    }
}

resource "vercel_project" "custom" {
    name = "test-acc-%[1]s-custom"
    %[2]s
}

resource "vercel_firewall_config" "custom" {
    project_id = vercel_project.custom.id
    %[2]s

    rules {
        rule {
          name =  "test"
          action = {
            action = "deny"
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "eq"
                value = "/test"
            },
            {
              type = "method"
              op = "eq"
              value = "POST"
            }]
          }]
        }
        rule {
          name =  "test2"
          action = {
            action = "rate_limit"
            rate_limit = {
                limit = 100
                window = 60
                algo = "fixed_window"
                keys = ["ip", "ja4"]
                action = "challenge"

            }
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "eq"
                value = "/api"
            }]
          }]
        }
        rule {
          name =  "test3"
          action = {
            action = "redirect"
            redirect = {
                location = "/bye"
                permanent = false

            }
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "eq"
                value = "/api"
            }]
          }]
        }

        rule {
          name =  "test4"
          action = {
            action = "log"
          }
          condition_group = [{
            conditions = [{
                type = "ja4_digest"
                op = "eq"
                value = "fakeja4"
            }]
          }]
        }
        rule {
          name =  "test_list"
          action = {
            action = "deny"
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "inc"
                values = [
                    "/test1",
                    "/test2",
                    "/test3"
                ]
            }]
          }]
        }
    }
}

resource "vercel_project" "ips" {
    name = "test-acc-%[1]s-ips"
    %[2]s
}

resource "vercel_firewall_config" "ips" {
    project_id = vercel_project.ips.id
    %[2]s

    ip_rules {
        rule {
            action = "deny"
            ip = "5.5.0.0/16"
            hostname = "test.com"
        }
        rule {
            action = "deny"
            ip = "1.2.3.4/32"
            hostname = "*.test.com"
        }
        rule {
            action = "deny"
            ip = "2.4.6.8"
            hostname = "*"
        }
    }
}

resource "vercel_project" "botfilter" {
    name = "test-acc-%[1]s-botfilter"
    %[2]s
}

resource "vercel_firewall_config" "botfilter" {
    project_id = vercel_project.botfilter.id
    %[2]s

    managed_rulesets {
        bot_filter {
            action = "challenge"
            active = true
        }
    }
}

`, name, teamID)
}

func testAccFirewallConfigResourceUpdated(name, teamID string) string {
	return fmt.Sprintf(`
resource "vercel_project" "managed" {
    name = "test-acc-%[1]s-mrs"
    %[2]s
}

resource "vercel_firewall_config" "managed" {
    project_id = vercel_project.managed.id
    %[2]s

    managed_rulesets {
        owasp {
            xss = { action = "deny", active = false }
            sqli = { action = "deny" }

            rce = { action = "deny", active = false }
            php = { action = "log" }
        }
    }
}

resource "vercel_project" "custom" {
    name = "test-acc-%[1]s-custom"
    %[2]s
}

resource "vercel_firewall_config" "custom" {
    project_id = vercel_project.custom.id
    %[2]s

    rules {
        rule {
          name =  "test1"
          action = {
            action = "deny"
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "eq"
                value = "/test"
            },
            {
              type = "method"
              op = "eq"
              value = "POST"
            }]
          }]
        }
        rule {
          name =  "test2"
          action = {
            action = "rate_limit"
            rate_limit = {
                limit = 150
                window = 60
                algo = "fixed_window"
                keys = ["ip"]
                action = "challenge"

            }
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "eq"
                value = "/api"
            }]
          }]
        }
        rule {
          name =  "test3"
          action = {
            action = "redirect"
            redirect = {
                location = "/bye"
                permanent = false

            }
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "eq"
                value = "/api"
            }]
          }]
        }
        rule {
          name =  "test_list"
          action = {
            action = "deny"
          }
          condition_group = [{
            conditions = [{
                type = "path"
                op = "inc"
                values = [
                    "/api",
                    "/api2",
                    "/api3"
                ]
            }]
          }]
        }
    }
}

resource "vercel_project" "ips" {
    name = "test-acc-%[1]s-ips"
    %[2]s
}

resource "vercel_firewall_config" "ips" {
    project_id = vercel_project.ips.id
    %[2]s

    ip_rules {
        rule {
            action = "deny"
            ip = "5.6.0.0/16"
            hostname = "test.com"
        }
        rule {
            action = "challenge"
            ip = "1.2.3.4/32"
            hostname = "*.test.com"
        }
        rule {
            action = "deny"
            ip = "2.4.6.8/32"
            hostname = "*"
        }
    }
}

resource "vercel_project" "neg" {
    name = "test-acc-%[1]s-neg"
    %[2]s
}

resource "vercel_firewall_config" "neg" {
    project_id = vercel_project.neg.id
    %[2]s

    rules {
        rule {
          name =  "test"
          action = {
            action = "deny"
          }
          condition_group = [{
            conditions = [{
                type = "ip_address"
                op = "inc"
                neg = true
                values = [
                    "1.2.3.4",
                    "3.4.5.6",
                    "5.6.7.7",
                ]
            }]
          }]
        }
    }
}

resource "vercel_project" "botfilter" {
    name = "test-acc-%[1]s-botfilter"
    %[2]s
}

resource "vercel_firewall_config" "botfilter" {
    project_id = vercel_project.botfilter.id
    %[2]s

    managed_rulesets {
        bot_filter {
            action = "deny"
            active = false
        }
    }
}
`, name, teamID)
}
