locals {
  enabled = true

  ip_sets = local.enabled && var.ip_set_reference_statement_rules != null ? {
    for indx, rule in flatten(var.ip_set_reference_statement_rules) :
    lookup(rule, "name", null) != null ? format("%s-ip-set", rule.name) : format("ip-set-%d", rule.priority)
    => rule.statement.ip_set if try(rule.statement.ip_set, null) != null && try(rule.statement.arn, null) == null
  } : {}

  ip_rule_to_ip_set = local.enabled && local.ip_set_reference_statement_rules != null ? {
    for name, rule in local.ip_set_reference_statement_rules :
    name => lookup(rule, "name", null) != null ? format("%s-ip-set", rule.name) : format("ip-set-%d", rule.priority)
  } : {}

  byte_match_statement_rules = local.enabled && var.byte_match_statement_rules != null ? {
    for rule in flatten(var.byte_match_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-byte-match-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  geo_allowlist_statement_rules = local.enabled && var.geo_allowlist_statement_rules != null ? {
    for rule in flatten(var.geo_allowlist_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-geo-allowlist-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  geo_match_statement_rules = local.enabled && var.geo_match_statement_rules != null ? {
    for rule in flatten(var.geo_match_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-geo-match-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  ip_set_reference_statement_rules = local.enabled && var.ip_set_reference_statement_rules != null ? {
    for indx, rule in flatten(var.ip_set_reference_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-ip-set-reference-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  managed_rule_group_statement_rules = local.enabled && var.managed_rule_group_statement_rules != null ? {
    for rule in flatten(var.managed_rule_group_statement_rules) :
    lookup(rule, "name", null) != null ? rule.name : format("%s-managed-rule-group-%d", module.this.id, rule.priority) => rule
  } : {}

  rate_based_statement_rules = local.enabled && var.rate_based_statement_rules != null ? {
    for rule in flatten(var.rate_based_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-rate-based-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  rule_group_reference_statement_rules = local.enabled && var.rule_group_reference_statement_rules != null ? {
    for rule in flatten(var.rule_group_reference_statement_rules) :
    lookup(rule, "name", null) != null ? rule.name : format("%s-rule-group-reference-%d", module.this.id, rule.priority) => rule
  } : {}

  regex_pattern_set_reference_statement_rules = local.enabled && var.regex_pattern_set_reference_statement_rules != null ? {
    for rule in flatten(var.regex_pattern_set_reference_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-regex-pattern-set-reference-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  regex_match_statement_rules = local.enabled && var.regex_match_statement_rules != null ? {
    for rule in flatten(var.regex_match_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-regex-match-statement-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  size_constraint_statement_rules = local.enabled && var.size_constraint_statement_rules != null ? {
    for rule in flatten(var.size_constraint_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-size-constraint-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  sqli_match_statement_rules = local.enabled && var.sqli_match_statement_rules != null ? {
    for rule in flatten(var.sqli_match_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-sqli-match-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  xss_match_statement_rules = local.enabled && var.xss_match_statement_rules != null ? {
    for rule in flatten(var.xss_match_statement_rules) :
    format("%s-%s",
      lookup(rule, "name", null) != null ? rule.name : format("%s-xss-match-%d", module.this.id, rule.priority),
      rule.action,
    ) => rule
  } : {}

  nested_rule_group_statement_rules = {
    for rule in var.nested_statement_rules != null ? var.nested_statement_rules : [] : rule.name => {
      name              = rule.name
      priority          = rule.priority
      action            = rule.action
      visibility_config = rule.visibility_config

      statements = [
        for stmt in rule.statement.and_statement.statements : {
          label_match_statement    = stmt.type == "label_match_statement" ? try(jsondecode(stmt.statement), null) : null
          not_byte_match_statement = stmt.type == "not_byte_match_statement" ? try(jsondecode(stmt.statement), null) : null
          # This is now easily extensible with other types:
          # not_geo_match_statement = stmt.type == "not_geo_match" ? jsondecode(stmt.statement) : null
        }
      ]
    }
  }

  default_custom_response_body_key = var.default_block_custom_response_body_key != null ? contains(keys(var.custom_response_body), var.default_block_custom_response_body_key) ? var.default_block_custom_response_body_key : null : null
}

# https://docs.aws.amazon.com/waf/latest/APIReference/API_AssociateWebACL.html
resource "aws_wafv2_web_acl_association" "default" {
  count = local.enabled && length(var.association_resource_arns) > 0 ? length(var.association_resource_arns) : 0

  resource_arn = var.association_resource_arns[count.index]
  web_acl_arn  = one(aws_wafv2_web_acl.default[*].arn)
}

# To start logging from a WAFv2 Web ACL, you need to create an Amazon Kinesis Data Firehose resource,
# such as the `aws_kinesis_firehose_delivery_stream` resource.
# Make sure to create the firehose with a PUT source (not a stream) in the region where you are operating.
# If you are capturing logs for Amazon CloudFront, create the firehose in the US East (N. Virginia) region.
# It is important to name the data firehose, CloudWatch log group, and/or S3 bucket with a prefix of `aws-waf-logs-`.
resource "aws_wafv2_web_acl_logging_configuration" "default" {
  count = local.enabled && length(var.log_destination_configs) > 0 ? 1 : 0

  resource_arn            = one(aws_wafv2_web_acl.default[*].arn)
  log_destination_configs = var.log_destination_configs

  dynamic "redacted_fields" {
    for_each = var.redacted_fields

    content {
      dynamic "method" {
        for_each = redacted_fields.value.method ? [true] : []
        content {}
      }

      dynamic "query_string" {
        for_each = redacted_fields.value.query_string ? [true] : []
        content {}
      }

      dynamic "uri_path" {
        for_each = redacted_fields.value.uri_path ? [true] : []
        content {}
      }

      dynamic "single_header" {
        for_each = lookup(redacted_fields.value, "single_header", null) != null ? toset(redacted_fields.value.single_header) : []
        content {
          name = single_header.value
        }
      }
    }
  }

  dynamic "logging_filter" {
    for_each = var.logging_filter != null ? [true] : []

    content {
      default_behavior = var.logging_filter.default_behavior

      dynamic "filter" {
        for_each = var.logging_filter.filter

        content {
          behavior    = filter.value.behavior
          requirement = filter.value.requirement

          dynamic "condition" {
            for_each = filter.value.condition

            content {
              dynamic "action_condition" {
                for_each = condition.value.action_condition != null ? [true] : []
                content {
                  action = condition.value.action_condition.action
                }
              }
              dynamic "label_name_condition" {
                for_each = condition.value.label_name_condition != null ? [true] : []
                content {
                  label_name = condition.value.label_name_condition.label_name
                }
              }
            }
          }
        }
      }
    }
  }
}

resource "aws_wafv2_web_acl" "default" {
  count = local.enabled ? 1 : 0

  name          = module.this.id
  description   = var.description
  scope         = var.scope
  token_domains = var.token_domains
  tags          = module.this.tags

  default_action {
    dynamic "allow" {
      for_each = var.default_action == "allow" ? [true] : []
      content {}
    }

    dynamic "block" {
      for_each = var.default_action == "block" ? [true] : []
      content {
        dynamic "custom_response" {
          for_each = var.default_block_response != null ? [true] : []
          content {
            response_code            = var.default_block_response
            custom_response_body_key = local.default_custom_response_body_key
          }
        }
      }
    }
  }

  visibility_config {
    cloudwatch_metrics_enabled = var.visibility_config.cloudwatch_metrics_enabled
    metric_name                = var.visibility_config.metric_name
    sampled_requests_enabled   = var.visibility_config.sampled_requests_enabled
  }

  dynamic "custom_response_body" {
    for_each = var.custom_response_body
    content {
      key          = custom_response_body.key
      content      = custom_response_body.value.content
      content_type = custom_response_body.value.content_type
    }
  }

  dynamic "rule" {
    for_each = local.byte_match_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []

          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []

          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []

          content {}
        }
      }

      statement {
        dynamic "byte_match_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            positional_constraint = byte_match_statement.value.positional_constraint
            search_string         = byte_match_statement.value.search_string

            dynamic "field_to_match" {
              for_each = lookup(rule.value.statement, "field_to_match", null) != null ? [rule.value.statement.field_to_match] : []

              content {
                dynamic "all_query_arguments" {
                  for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                  content {}
                }

                dynamic "body" {
                  for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                  content {}
                }

                dynamic "method" {
                  for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                  content {}
                }

                dynamic "query_string" {
                  for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                  content {}
                }

                dynamic "single_header" {
                  for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                  content {
                    name = single_header.value.name
                  }
                }

                dynamic "single_query_argument" {
                  for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                  content {
                    name = single_query_argument.value.name
                  }
                }

                dynamic "uri_path" {
                  for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                  content {}
                }
              }
            }

            dynamic "text_transformation" {
              for_each = lookup(rule.value.statement, "text_transformation", null) != null ? [
                for rule in lookup(rule.value.statement, "text_transformation") : {
                  priority = rule.priority
                  type     = rule.type
              }] : []

              content {
                priority = text_transformation.value.priority
                type     = text_transformation.value.type
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.geo_allowlist_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
      }

      # `geo_allowlist_statement_rules` is a special case where we use `not_statement` to wrap our `statement` block to support
      # an "allowlist". Otherwise, using `geo_match_statement_rules` requires specifying ALL country codes that you
      # would like to blocklist.
      statement {
        not_statement {
          statement {
            dynamic "geo_match_statement" {
              for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

              content {
                country_codes = geo_match_statement.value.country_codes

                dynamic "forwarded_ip_config" {
                  for_each = lookup(geo_match_statement.value, "forwarded_ip_config", null) != null ? [geo_match_statement.value.forwarded_ip_config] : []

                  content {
                    fallback_behavior = forwarded_ip_config.value.fallback_behavior
                    header_name       = forwarded_ip_config.value.header_name
                  }
                }
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.geo_match_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "geo_match_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            country_codes = geo_match_statement.value.country_codes

            dynamic "forwarded_ip_config" {
              for_each = lookup(geo_match_statement.value, "forwarded_ip_config", null) != null ? [geo_match_statement.value.forwarded_ip_config] : []

              content {
                fallback_behavior = forwarded_ip_config.value.fallback_behavior
                header_name       = forwarded_ip_config.value.header_name
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.ip_set_reference_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "ip_set_reference_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            arn = try(aws_wafv2_ip_set.default[local.ip_rule_to_ip_set[rule.key]], null) != null ? aws_wafv2_ip_set.default[local.ip_rule_to_ip_set[rule.key]].arn : ip_set_reference_statement.value.arn

            dynamic "ip_set_forwarded_ip_config" {
              for_each = lookup(ip_set_reference_statement.value, "ip_set_forwarded_ip_config", null) != null ? [ip_set_reference_statement.value.ip_set_forwarded_ip_config] : []

              content {
                fallback_behavior = ip_set_forwarded_ip_config.value.fallback_behavior
                header_name       = ip_set_forwarded_ip_config.value.header_name
                position          = ip_set_forwarded_ip_config.value.position
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.managed_rule_group_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      override_action {
        dynamic "count" {
          for_each = lookup(rule.value, "override_action", null) == "count" ? [1] : []
          content {}
        }
        dynamic "none" {
          for_each = lookup(rule.value, "override_action", null) != "count" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "managed_rule_group_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            name        = managed_rule_group_statement.value.name
            vendor_name = managed_rule_group_statement.value.vendor_name
            version     = lookup(managed_rule_group_statement.value, "version", null)

            dynamic "rule_action_override" {
              for_each = lookup(managed_rule_group_statement.value, "rule_action_override", null) != null ? managed_rule_group_statement.value.rule_action_override : {}

              content {
                name = rule_action_override.key

                # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#action-block
                action_to_use {
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#allow-block
                  dynamic "allow" {
                    for_each = rule_action_override.value.action == "allow" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#block-block
                  dynamic "block" {
                    for_each = rule_action_override.value.action == "block" ? [1] : []
                    content {
                      dynamic "custom_response" {
                        for_each = lookup(rule_action_override.value, "custom_response", null) != null ? [1] : []
                        content {
                          response_code            = rule_action_override.value.custom_response.response_code
                          custom_response_body_key = lookup(rule_action_override.value.custom_response, "custom_response_body_key", null)
                          dynamic "response_header" {
                            for_each = lookup(rule_action_override.value.custom_response, "response_header", null) != null ? [1] : []
                            content {
                              name  = rule_action_override.value.custom_response.response_header.name
                              value = rule_action_override.value.custom_response.response_header.value
                            }
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#count-block
                  dynamic "count" {
                    for_each = rule_action_override.value.action == "count" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#captcha-block
                  dynamic "captcha" {
                    for_each = rule_action_override.value.action == "captcha" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#challenge-block
                  dynamic "challenge" {
                    for_each = rule_action_override.value.action == "challenge" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                }
              }
            }

            # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#managed_rule_group_configs-block
            dynamic "managed_rule_group_configs" {
              for_each = lookup(managed_rule_group_statement.value, "managed_rule_group_configs", null) != null ? managed_rule_group_statement.value.managed_rule_group_configs : []
              content {
                dynamic "aws_managed_rules_anti_ddos_rule_set" {
                  for_each = lookup(managed_rule_group_configs.value, "aws_managed_rules_anti_ddos_rule_set", null) != null ? [1] : []
                  content {
                    sensitivity_to_block = managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set.sensitivity_to_block
                    dynamic "client_side_action_config" {
                      for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set, "client_side_action_config", null) != null ? [managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set.client_side_action_config] : []
                      content {
                        challenge {
                          usage_of_action = managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set.client_side_action_config.challenge.usage_of_action
                          sensitivity     = lookup(managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set.client_side_action_config.challenge, "sensitivity", null)
                          dynamic "exempt_uri_regular_expression" {
                            for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set.client_side_action_config.challenge, "exempt_uri_regular_expression", null) != null ? managed_rule_group_configs.value.aws_managed_rules_anti_ddos_rule_set.client_side_action_config.challenge.exempt_uri_regular_expression : []
                            content {
                              regex_string = exempt_uri_regular_expression.value.regex_string
                            }
                          }
                        }
                      }
                    }
                  }
                }

                dynamic "aws_managed_rules_bot_control_rule_set" {
                  for_each = lookup(managed_rule_group_configs.value, "aws_managed_rules_bot_control_rule_set", null) != null ? [1] : []
                  content {
                    inspection_level        = managed_rule_group_configs.value.aws_managed_rules_bot_control_rule_set.inspection_level
                    enable_machine_learning = managed_rule_group_configs.value.aws_managed_rules_bot_control_rule_set.enable_machine_learning
                  }
                }

                dynamic "aws_managed_rules_atp_rule_set" {
                  for_each = lookup(managed_rule_group_configs.value, "aws_managed_rules_atp_rule_set", null) != null ? [1] : []
                  content {
                    enable_regex_in_path = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set, "enable_regex_in_path", null)
                    login_path           = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.login_path

                    dynamic "request_inspection" {
                      for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set, "request_inspection", null) != null ? [1] : []
                      content {
                        payload_type = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.request_inspection.payload_type
                        username_field {
                          identifier = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.request_inspection.username_field.identifier
                        }
                        password_field {
                          identifier = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.request_inspection.password_field.identifier
                        }
                      }
                    }

                    dynamic "response_inspection" {
                      for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set, "response_inspection", null) != null ? [1] : []
                      content {
                        dynamic "body_contains" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection, "body_contains", null) != null ? [1] : []
                          content {
                            failure_strings = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.body_contains.failure_strings
                            success_strings = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.body_contains.success_strings
                          }
                        }
                        dynamic "header" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection, "header", null) != null ? [1] : []
                          content {
                            failure_values = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.header.failure_values
                            name           = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.header.name
                            success_values = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.header.success_values
                          }
                        }
                        dynamic "json" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection, "json", null) != null ? [1] : []
                          content {
                            failure_values = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.json.failure_values
                            identifier     = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.json.identifier
                            success_values = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.json.success_values
                          }
                        }
                        dynamic "status_code" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection, "status_code", null) != null ? [1] : []
                          content {
                            failure_codes = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.status_code.failure_codes
                            success_codes = managed_rule_group_configs.value.aws_managed_rules_atp_rule_set.response_inspection.status_code.success_codes
                          }
                        }
                      }
                    }
                  }
                }

                dynamic "aws_managed_rules_acfp_rule_set" {
                  for_each = lookup(managed_rule_group_configs.value, "aws_managed_rules_acfp_rule_set", null) != null ? [1] : []
                  content {
                    creation_path          = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.creation_path
                    enable_regex_in_path   = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set, "enable_regex_in_path", true)
                    registration_page_path = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.registration_page_path

                    dynamic "request_inspection" {
                      for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set, "request_inspection", null) != null ? [1] : []
                      content {
                        payload_type = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection.payload_type
                        dynamic "username_field" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection, "username_field", null) != null ? [1] : []
                          content {
                            identifier = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection.username_field.identifier
                          }
                        }
                        dynamic "password_field" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection, "password_field", null) != null ? [1] : []
                          content {
                            identifier = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection.password_field.identifier
                          }
                        }
                        dynamic "email_field" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection, "email_field", null) != null ? [1] : []
                          content {
                            identifier = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection.email_field.identifier
                          }
                        }
                        dynamic "address_fields" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection, "address_fields", null) != null ? [1] : []
                          content {
                            identifiers = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection.address_fields.identifiers
                          }
                        }
                        dynamic "phone_number_fields" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection, "phone_number_fields", null) != null ? [1] : []
                          content {
                            identifiers = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.request_inspection.phone_number_fields.identifiers
                          }
                        }
                      }
                    }

                    dynamic "response_inspection" {
                      for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set, "response_inspection", null) != null ? [1] : []
                      content {
                        dynamic "body_contains" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection, "body_contains", null) != null ? [1] : []
                          content {
                            failure_strings = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.body_contains.failure_strings
                            success_strings = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.body_contains.success_strings
                          }
                        }
                        dynamic "header" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection, "header", null) != null ? [1] : []
                          content {
                            failure_values = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.header.failure_values
                            name           = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.header.name
                            success_values = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.header.success_values
                          }
                        }
                        dynamic "json" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection, "json", null) != null ? [1] : []
                          content {
                            failure_values = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.json.failure_values
                            identifier     = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.json.identifier
                            success_values = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.json.success_values
                          }
                        }
                        dynamic "status_code" {
                          for_each = lookup(managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection, "status_code", null) != null ? [1] : []
                          content {
                            failure_codes = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.status_code.failure_codes
                            success_codes = managed_rule_group_configs.value.aws_managed_rules_acfp_rule_set.response_inspection.status_code.success_codes
                          }
                        }
                      }
                    }
                  }
                }
              }
            }

            dynamic "scope_down_statement" {
              for_each = lookup(managed_rule_group_statement.value, "scope_down_statement", null) != null && !lookup(managed_rule_group_statement.value, "scope_down_not_statement_enabled", false) ? [managed_rule_group_statement.value.scope_down_statement] : []

              content {
                dynamic "byte_match_statement" {
                  for_each = lookup(scope_down_statement.value, "byte_match_statement", null) != null ? [scope_down_statement.value.byte_match_statement] : []

                  content {
                    positional_constraint = byte_match_statement.value.positional_constraint
                    search_string         = byte_match_statement.value.search_string

                    dynamic "field_to_match" {
                      for_each = lookup(byte_match_statement.value, "field_to_match", null) != null ? [byte_match_statement.value.field_to_match] : []

                      content {
                        dynamic "all_query_arguments" {
                          for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []
                          content {}
                        }

                        dynamic "body" {
                          for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []
                          content {}
                        }

                        dynamic "method" {
                          for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []
                          content {}
                        }

                        dynamic "query_string" {
                          for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []
                          content {}
                        }

                        dynamic "single_header" {
                          for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []
                          content {
                            name = single_header.value.name
                          }
                        }

                        dynamic "single_query_argument" {
                          for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []
                          content {
                            name = single_query_argument.value.name
                          }
                        }

                        dynamic "uri_path" {
                          for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []
                          content {}
                        }
                      }
                    }

                    dynamic "text_transformation" {
                      for_each = lookup(byte_match_statement.value, "text_transformation", null) != null ? [
                        for rule in byte_match_statement.value.text_transformation : {
                          priority = rule.priority
                          type     = rule.type
                      }] : []

                      content {
                        priority = text_transformation.value.priority
                        type     = text_transformation.value.type
                      }
                    }
                  }
                }
              }
            }

            dynamic "scope_down_statement" {
              for_each = lookup(managed_rule_group_statement.value, "scope_down_statement", null) != null && lookup(managed_rule_group_statement.value, "scope_down_not_statement_enabled", false) ? [managed_rule_group_statement.value.scope_down_statement] : []

              content {
                not_statement {
                  statement {
                    dynamic "byte_match_statement" {
                      for_each = lookup(scope_down_statement.value, "byte_match_statement", null) != null ? [scope_down_statement.value.byte_match_statement] : []

                      content {
                        positional_constraint = byte_match_statement.value.positional_constraint
                        search_string         = byte_match_statement.value.search_string

                        dynamic "field_to_match" {
                          for_each = lookup(byte_match_statement.value, "field_to_match", null) != null ? [byte_match_statement.value.field_to_match] : []

                          content {
                            dynamic "all_query_arguments" {
                              for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []
                              content {}
                            }

                            dynamic "body" {
                              for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []
                              content {}
                            }

                            dynamic "method" {
                              for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []
                              content {}
                            }

                            dynamic "query_string" {
                              for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []
                              content {}
                            }

                            dynamic "single_header" {
                              for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []
                              content {
                                name = single_header.value.name
                              }
                            }

                            dynamic "single_query_argument" {
                              for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []
                              content {
                                name = single_query_argument.value.name
                              }
                            }

                            dynamic "uri_path" {
                              for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []
                              content {}
                            }
                          }
                        }

                        dynamic "text_transformation" {
                          for_each = lookup(byte_match_statement.value, "text_transformation", null) != null ? [
                            for rule in byte_match_statement.value.text_transformation : {
                              priority = rule.priority
                              type     = rule.type
                          }] : []

                          content {
                            priority = text_transformation.value.priority
                            type     = text_transformation.value.type
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.rate_based_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "rate_based_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            aggregate_key_type    = lookup(rate_based_statement.value, "aggregate_key_type", "IP")
            limit                 = rate_based_statement.value.limit
            evaluation_window_sec = lookup(rate_based_statement.value, "evaluation_window_sec", 300)

            dynamic "forwarded_ip_config" {
              for_each = lookup(rate_based_statement.value, "forwarded_ip_config", null) != null ? [rate_based_statement.value.forwarded_ip_config] : []

              content {
                fallback_behavior = forwarded_ip_config.value.fallback_behavior
                header_name       = forwarded_ip_config.value.header_name
              }
            }

            dynamic "scope_down_statement" {
              for_each = lookup(rate_based_statement.value, "scope_down_statement", null) != null ? [rate_based_statement.value.scope_down_statement] : []

              content {
                dynamic "byte_match_statement" {
                  for_each = lookup(scope_down_statement.value, "byte_match_statement", null) != null ? [scope_down_statement.value.byte_match_statement] : []

                  content {
                    positional_constraint = byte_match_statement.value.positional_constraint
                    search_string         = byte_match_statement.value.search_string

                    dynamic "field_to_match" {
                      for_each = lookup(byte_match_statement.value, "field_to_match", null) != null ? [byte_match_statement.value.field_to_match] : []

                      content {
                        dynamic "all_query_arguments" {
                          for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                          content {}
                        }

                        dynamic "body" {
                          for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                          content {}
                        }

                        dynamic "method" {
                          for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                          content {}
                        }

                        dynamic "query_string" {
                          for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                          content {}
                        }

                        dynamic "single_header" {
                          for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                          content {
                            name = single_header.value.name
                          }
                        }

                        dynamic "single_query_argument" {
                          for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                          content {
                            name = single_query_argument.value.name
                          }
                        }

                        dynamic "uri_path" {
                          for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                          content {}
                        }
                      }
                    }

                    dynamic "text_transformation" {
                      for_each = lookup(byte_match_statement.value, "text_transformation", null) != null ? [
                        for rule in byte_match_statement.value.text_transformation : {
                          priority = rule.priority
                          type     = rule.type
                      }] : []

                      content {
                        priority = text_transformation.value.priority
                        type     = text_transformation.value.type
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.regex_pattern_set_reference_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "regex_pattern_set_reference_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            arn = regex_pattern_set_reference_statement.value.arn

            dynamic "field_to_match" {
              for_each = lookup(rule.value.statement, "field_to_match", null) != null ? [rule.value.statement.field_to_match] : []

              content {
                dynamic "all_query_arguments" {
                  for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                  content {}
                }

                dynamic "body" {
                  for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                  content {}
                }

                dynamic "method" {
                  for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                  content {}
                }

                dynamic "query_string" {
                  for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                  content {}
                }

                dynamic "single_header" {
                  for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                  content {
                    name = single_header.value.name
                  }
                }

                dynamic "single_query_argument" {
                  for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                  content {
                    name = single_query_argument.value.name
                  }
                }

                dynamic "uri_path" {
                  for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                  content {}
                }
              }
            }

            dynamic "text_transformation" {
              for_each = lookup(rule.value.statement, "text_transformation", null) != null ? [
                for rule in lookup(rule.value.statement, "text_transformation") : {
                  priority = rule.priority
                  type     = rule.type
              }] : []

              content {
                priority = text_transformation.value.priority
                type     = text_transformation.value.type
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.regex_match_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "regex_match_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            regex_string = regex_match_statement.value.regex_string

            dynamic "field_to_match" {
              for_each = lookup(rule.value.statement, "field_to_match", null) != null ? [rule.value.statement.field_to_match] : []

              content {
                dynamic "all_query_arguments" {
                  for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                  content {}
                }

                dynamic "body" {
                  for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                  content {}
                }

                dynamic "method" {
                  for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                  content {}
                }

                dynamic "query_string" {
                  for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                  content {}
                }

                dynamic "single_header" {
                  for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                  content {
                    name = single_header.value.name
                  }
                }

                dynamic "single_query_argument" {
                  for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                  content {
                    name = single_query_argument.value.name
                  }
                }

                dynamic "uri_path" {
                  for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                  content {}
                }
              }
            }

            dynamic "text_transformation" {
              for_each = lookup(rule.value.statement, "text_transformation", null) != null ? [
                for rule in lookup(rule.value.statement, "text_transformation") : {
                  priority = rule.priority
                  type     = rule.type
              }] : []

              content {
                priority = text_transformation.value.priority
                type     = text_transformation.value.type
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.rule_group_reference_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      override_action {
        dynamic "count" {
          for_each = lookup(rule.value, "override_action", null) == "count" ? [1] : []
          content {}
        }
        dynamic "none" {
          for_each = lookup(rule.value, "override_action", null) != "count" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "rule_group_reference_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            arn = rule_group_reference_statement.value.arn

            dynamic "rule_action_override" {
              for_each = lookup(rule_group_reference_statement.value, "rule_action_override", null) != null ? rule_group_reference_statement.value.rule_action_override : {}

              content {
                name = rule_action_override.key

                # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#action-block
                action_to_use {
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#allow-block
                  dynamic "allow" {
                    for_each = rule_action_override.value.action == "allow" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#block-block
                  dynamic "block" {
                    for_each = rule_action_override.value.action == "block" ? [1] : []
                    content {
                      dynamic "custom_response" {
                        for_each = lookup(rule_action_override.value, "custom_response", null) != null ? [1] : []
                        content {
                          response_code            = rule_action_override.value.custom_response.response_code
                          custom_response_body_key = lookup(rule_action_override.value.custom_response, "custom_response_body_key", null)
                          dynamic "response_header" {
                            for_each = lookup(rule_action_override.value.custom_response, "response_header", null) != null ? [1] : []
                            content {
                              name  = rule_action_override.value.custom_response.response_header.name
                              value = rule_action_override.value.custom_response.response_header.value
                            }
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#count-block
                  dynamic "count" {
                    for_each = rule_action_override.value.action == "count" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#captcha-block
                  dynamic "captcha" {
                    for_each = rule_action_override.value.action == "captcha" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/wafv2_web_acl#challenge-block
                  dynamic "challenge" {
                    for_each = rule_action_override.value.action == "challenge" ? [1] : []
                    content {
                      dynamic "custom_request_handling" {
                        for_each = lookup(rule_action_override.value, "custom_request_handling", null) != null ? [1] : []
                        content {
                          insert_header {
                            name  = rule_action_override.value.custom_request_handling.insert_header.name
                            value = rule_action_override.value.custom_request_handling.insert_header.value
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.size_constraint_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "size_constraint_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {
            comparison_operator = size_constraint_statement.value.comparison_operator
            size                = size_constraint_statement.value.size

            dynamic "field_to_match" {
              for_each = lookup(rule.value.statement, "field_to_match", null) != null ? [rule.value.statement.field_to_match] : []

              content {
                dynamic "all_query_arguments" {
                  for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                  content {}
                }

                dynamic "body" {
                  for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                  content {
                    # Oversize handling tells AWS WAF what to do with a web request when the request component that the rule inspects is over the limits.
                    # WAF does not support inspecting the entire contents of the body of a web request when the body exceeds 8 KB (8192 bytes).
                    # Only the first 8 KB of the request body are forwarded to WAF by the underlying host service
                    # Valid values include the following: CONTINUE, MATCH, NO_MATCH
                    oversize_handling = try(field_to_match.value.body.oversize_handling, "CONTINUE")
                  }
                }

                dynamic "method" {
                  for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                  content {}
                }

                dynamic "query_string" {
                  for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                  content {}
                }

                dynamic "single_header" {
                  for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                  content {
                    name = single_header.value.name
                  }
                }

                dynamic "single_query_argument" {
                  for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                  content {
                    name = single_query_argument.value.name
                  }
                }

                dynamic "uri_path" {
                  for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                  content {}
                }
              }
            }

            dynamic "text_transformation" {
              for_each = lookup(rule.value.statement, "text_transformation", null) != null ? [
                for rule in lookup(rule.value.statement, "text_transformation") : {
                  priority = rule.priority
                  type     = rule.type
              }] : []

              content {
                priority = text_transformation.value.priority
                type     = text_transformation.value.type
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.sqli_match_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "sqli_match_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {

            dynamic "field_to_match" {
              for_each = lookup(rule.value.statement, "field_to_match", null) != null ? [rule.value.statement.field_to_match] : []

              content {
                dynamic "all_query_arguments" {
                  for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                  content {}
                }

                dynamic "body" {
                  for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                  content {}
                }

                dynamic "method" {
                  for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                  content {}
                }

                dynamic "query_string" {
                  for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                  content {}
                }

                dynamic "single_header" {
                  for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                  content {
                    name = single_header.value.name
                  }
                }

                dynamic "single_query_argument" {
                  for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                  content {
                    name = single_query_argument.value.name
                  }
                }

                dynamic "uri_path" {
                  for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                  content {}
                }
              }
            }

            dynamic "text_transformation" {
              for_each = lookup(rule.value.statement, "text_transformation", null) != null ? [
                for rule in lookup(rule.value.statement, "text_transformation") : {
                  priority = rule.priority
                  type     = rule.type
              }] : []

              content {
                priority = text_transformation.value.priority
                type     = text_transformation.value.type
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.xss_match_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        dynamic "xss_match_statement" {
          for_each = lookup(rule.value, "statement", null) != null ? [rule.value.statement] : []

          content {

            dynamic "field_to_match" {
              for_each = lookup(rule.value.statement, "field_to_match", null) != null ? [rule.value.statement.field_to_match] : []

              content {
                dynamic "all_query_arguments" {
                  for_each = lookup(field_to_match.value, "all_query_arguments", null) != null ? [1] : []

                  content {}
                }

                dynamic "body" {
                  for_each = lookup(field_to_match.value, "body", null) != null ? [1] : []

                  content {}
                }

                dynamic "method" {
                  for_each = lookup(field_to_match.value, "method", null) != null ? [1] : []

                  content {}
                }

                dynamic "query_string" {
                  for_each = lookup(field_to_match.value, "query_string", null) != null ? [1] : []

                  content {}
                }

                dynamic "single_header" {
                  for_each = lookup(field_to_match.value, "single_header", null) != null ? [field_to_match.value.single_header] : []

                  content {
                    name = single_header.value.name
                  }
                }

                dynamic "single_query_argument" {
                  for_each = lookup(field_to_match.value, "single_query_argument", null) != null ? [field_to_match.value.single_query_argument] : []

                  content {
                    name = single_query_argument.value.name
                  }
                }

                dynamic "uri_path" {
                  for_each = lookup(field_to_match.value, "uri_path", null) != null ? [1] : []

                  content {}
                }
              }
            }

            dynamic "text_transformation" {
              for_each = lookup(rule.value.statement, "text_transformation", null) != null ? [
                for rule in lookup(rule.value.statement, "text_transformation") : {
                  priority = rule.priority
                  type     = rule.type
              }] : []

              content {
                priority = text_transformation.value.priority
                type     = text_transformation.value.type
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }

  dynamic "rule" {
    for_each = local.nested_rule_group_statement_rules

    content {
      name     = rule.value.name
      priority = rule.value.priority

      action {
        dynamic "allow" {
          for_each = rule.value.action == "allow" ? [1] : []
          content {}
        }
        dynamic "block" {
          for_each = rule.value.action == "block" ? [1] : []
          content {}
        }
        dynamic "count" {
          for_each = rule.value.action == "count" ? [1] : []
          content {}
        }
        dynamic "captcha" {
          for_each = rule.value.action == "captcha" ? [1] : []
          content {}
        }
      }

      statement {
        and_statement {
          dynamic "statement" {
            for_each = rule.value.statements
            content {
              dynamic "label_match_statement" {
                for_each = statement.value.label_match_statement != null ? [statement.value.label_match_statement] : []
                content {
                  scope = label_match_statement.value.scope
                  key   = label_match_statement.value.key
                }
              }
              dynamic "not_statement" {
                for_each = statement.value.not_byte_match_statement != null ? [1] : []
                content {
                  statement {
                    dynamic "byte_match_statement" {
                      for_each = statement.value.not_byte_match_statement != null ? [statement.value.not_byte_match_statement] : []
                      content {
                        positional_constraint = byte_match_statement.value.positional_constraint
                        search_string         = byte_match_statement.value.search_string

                        field_to_match {
                          dynamic "uri_path" {
                            for_each = try(byte_match_statement.value.field_to_match.uri_path, null) != null ? [{}] : []
                            content {}
                          }

                          dynamic "single_header" {
                            for_each = try(byte_match_statement.value.field_to_match.single_header, null) != null ? [byte_match_statement.value.field_to_match.single_header] : []
                            content {
                              name = single_header.value.name
                            }
                          }
                        }

                        dynamic "text_transformation" {
                          for_each = byte_match_statement.value.text_transformation
                          content {
                            priority = text_transformation.value.priority
                            type     = text_transformation.value.type
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }

      dynamic "visibility_config" {
        for_each = lookup(rule.value, "visibility_config", null) != null ? [rule.value.visibility_config] : []

        content {
          cloudwatch_metrics_enabled = lookup(visibility_config.value, "cloudwatch_metrics_enabled", true)
          metric_name                = visibility_config.value.metric_name
          sampled_requests_enabled   = lookup(visibility_config.value, "sampled_requests_enabled", true)
        }
      }

      dynamic "captcha_config" {
        for_each = lookup(rule.value, "captcha_config", null) != null ? [rule.value.captcha_config] : []

        content {
          immunity_time_property {
            immunity_time = captcha_config.value.immunity_time_property.immunity_time
          }
        }
      }

      dynamic "rule_label" {
        for_each = lookup(rule.value, "rule_label", null) != null ? rule.value.rule_label : []
        content {
          name = rule_label.value
        }
      }
    }
  }
}

module "ip_set_label" {
  for_each = local.ip_sets

  source  = "cloudposse/label/null"
  version = "0.25.0"

  attributes = [each.key]
  context    = module.this.context
}

resource "aws_wafv2_ip_set" "default" {
  for_each = local.ip_sets

  name               = module.ip_set_label[each.key].id
  description        = lookup(each.value, "description", null)
  scope              = var.scope
  ip_address_version = each.value.ip_address_version
  addresses          = each.value.addresses

  tags = module.this.tags
}
