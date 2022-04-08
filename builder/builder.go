package builder

import (
	"fmt"
	"strconv"
	"time"

	bsonMGO "github.com/biges/mgo/bson"
)

// Fields returns db filter query.
func Fields(filter map[string][]string) (bsonMGO.M, error) {
	dbQuery := bsonMGO.M{
		"deleted_at": nil,
	}

	// date filters
	fromCreatedAt, fromCheck := filter["from_created_at"]
	toCreatedAt, toCheck := filter["to_created_at"]
	if fromCheck || toCheck {
		var createdAtQuery = bsonMGO.M{}
		if fromCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", fromCreatedAt[0])
			if err != nil {
				return nil, err
			}
			createdAtQuery["$gte"] = t
		}

		if toCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", toCreatedAt[0])
			if err != nil {
				return nil, err
			}
			createdAtQuery["$lt"] = t
		}

		dbQuery["created_at"] = createdAtQuery
	}

	fromUpdatedAt, fromCheck := filter["from_updated_at"]
	toUpdatedAt, toCheck := filter["to_updated_at"]
	if fromCheck || toCheck {
		var updatedAtQuery = bsonMGO.M{}
		if fromCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", fromUpdatedAt[0])
			if err != nil {
				return nil, err
			}
			updatedAtQuery["$gte"] = t
		}

		if toCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", toUpdatedAt[0])
			if err != nil {
				return nil, err
			}
			updatedAtQuery["$lt"] = t
		}

		dbQuery["updated_at"] = updatedAtQuery
	}

	fromLastSeen, fromCheck := filter["from_last_seen"]
	toLastSeen, toCheck := filter["to_last_seen"]
	if fromCheck || toCheck {
		var updatedAtQuery = bsonMGO.M{}
		if fromCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", fromLastSeen[0])
			if err != nil {
				return nil, err
			}
			updatedAtQuery["$gte"] = t
		}

		if toCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", toLastSeen[0])
			if err != nil {
				return nil, err
			}
			updatedAtQuery["$lt"] = t
		}

		dbQuery["last_seen"] = updatedAtQuery
	}

	fromLastStatusSignal, fromCheck := filter["from_last_status"]
	toLastStatusSignal, toCheck := filter["to_last_status"]
	if fromCheck || toCheck {
		var lastStatusAtQuery = bsonMGO.M{}
		if fromCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", fromLastStatusSignal[0])
			if err != nil {
				return nil, err
			}
			lastStatusAtQuery["$gte"] = t
		}

		if toCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", toLastStatusSignal[0])
			if err != nil {
				return nil, err
			}
			lastStatusAtQuery["$lt"] = t
		}

		dbQuery["last_status_signal.created_at"] = lastStatusAtQuery
	}

	fromApprovedAt, fromCheck := filter["from_approved_at"]
	toApprovedAt, toCheck := filter["to_approved_at"]
	if fromCheck || toCheck {
		var approveAtQuery = bsonMGO.M{}
		if fromCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", fromApprovedAt[0])
			if err != nil {
				return nil, err
			}
			approveAtQuery["$gte"] = t
		}

		if toCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", toApprovedAt[0])
			if err != nil {
				return nil, err
			}
			approveAtQuery["$lt"] = t
		}

		dbQuery["approved_at"] = approveAtQuery
	}

	visitStartAt, fromCheck := filter["visit_start_at"]
	visitEndAt, toCheck := filter["visit_end_at"]
	if fromCheck || toCheck {
		if fromCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", visitStartAt[0])
			if err != nil {
				return nil, err
			}

			dbQuery["visit_start_at"] = bsonMGO.M{
				"$gte": t,
			}
		}

		if toCheck {
			t, err := time.Parse("2006-01-02T15:04:05.000Z", visitEndAt[0])
			if err != nil {
				return nil, err
			}

			dbQuery["visit_end_at"] = bsonMGO.M{
				"$lt": t,
			}
		}
	}

	// directly id filters
	if val, ok := filter["ids"]; ok {
		dbQuery["token"] = bsonMGO.M{
			"$in": val,
		}
	}

	if val, ok := filter["distributor_id"]; ok {
		dbQuery["distributor_id"] = val[0]
	}

	if val, ok := filter["reseller_id"]; ok {
		dbQuery["reseller_id"] = val[0]
	}

	if val, ok := filter["reseller_ids"]; ok {
		dbQuery["reseller_id"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["account_id"]; ok {
		dbQuery["account_id"] = val[0]
	}

	if val, ok := filter["premise_id"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["premise.token"] = val[0]
		} else {
			dbQuery["premise_id"] = val[0]
		}
	}

	if val, ok := filter["premise_ids"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["premise.token"] = bsonMGO.M{"$in": val}
		} else {
			dbQuery["premise_id"] = bsonMGO.M{"$in": val}
		}
	}

	if val, ok := filter["area_id"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["area.token"] = val[0]
		} else {
			dbQuery["area_id"] = val[0]
		}
	}

	if val, ok := filter["user_id"]; ok {
		dbQuery["user_id"] = val[0]
	}

	if val, ok := filter["device_id"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["device.token"] = val[0]
		} else {
			dbQuery["device_id"] = val[0]
		}
	}

	if val, ok := filter["device_ids"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["device.token"] = bsonMGO.M{"$in": val}
		} else {
			dbQuery["device_id"] = bsonMGO.M{"$in": val}
		}
	}

	if val, ok := filter["devices_ids"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["devices.token"] = bsonMGO.M{"$in": val}
		} else {
			dbQuery["device_ids"] = bsonMGO.M{"$in": val}
		}
	}

	if val, ok := filter["joined_accounts_ids"]; ok {
		dbQuery["accounts.account_id"] = bsonMGO.M{
			"$in": val,
		}
	}

	if val, ok := filter["notification_preferences_signal_categories"]; ok {
		dbQuery["notification_preferences.push_notification"] = bsonMGO.M{
			"$in": val,
		}
	}

	if val, ok := filter["joined_account_id"]; ok {
		dbQuery["accounts.account_id"] = val[0]
	}

	if val, ok := filter["joined_account_user_group_id"]; ok {
		dbQuery["accounts.user_group_id"] = val[0]
	}

	if val, ok := filter["subscription_id"]; ok {
		dbQuery["subscription_id"] = val[0]
	}

	if val, ok := filter["assignee_id"]; ok {
		dbQuery["assignee_id"] = val[0]
	}

	if val, ok := filter["key"]; ok {
		dbQuery["key"] = val[0]
	}

	// token filters
	if val, ok := filter["client_devices_check"]; ok {
		dbQuery[fmt.Sprintf("client_devices.%s.jwt_token", val[0])] = val[1]
	}

	// address filters
	if val, ok := filter["postal_code"]; ok {
		dbQuery["address.postal_code"] = val[0]
	}

	if val, ok := filter["country"]; ok {
		dbQuery["address.country"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["city"]; ok {
		dbQuery["address.city"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["province"]; ok {
		dbQuery["address.province"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["district"]; ok {
		dbQuery["address.district"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["street"]; ok {
		dbQuery["address.street"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["job_type"]; ok {
		dbQuery["job.type"] = val[0]
	}

	if val, ok := filter["job_distributor_id"]; ok {
		dbQuery["job.distributor_id"] = val[0]
	}

	if val, ok := filter["job_reseller_id"]; ok {
		dbQuery["job.reseller_id"] = val[0]
	}

	if val, ok := filter["job_user_group_id"]; ok {
		dbQuery["job.user_group_id"] = val[0]
	}

	if val, ok := filter["reset_password_token"]; ok {
		dbQuery["reset_password.token"] = val[0]
	}

	if val, ok := filter["action"]; ok {
		dbQuery["action"] = val[0]
	}

	// regex filters
	if val, ok := filter["email"]; ok {
		dbQuery["email"] = val[0]
	}

	if val, ok := filter["email_reg"]; ok {
		dbQuery["email"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["email_verified"]; ok {
		b, err := strconv.ParseBool(val[0])
		if err != nil {
			return nil, err
		}
		dbQuery["email_verified"] = b
	}

	if val, ok := filter["phone_number"]; ok {
		dbQuery["gsm"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["phone_number_verified"]; ok {
		b, err := strconv.ParseBool(val[0])
		if err != nil {
			return nil, err
		}
		dbQuery["phone_number_verified"] = b
	}

	if val, ok := filter["first_name"]; ok {
		dbQuery["first_name"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["last_name"]; ok {
		dbQuery["last_name"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["name"]; ok {
		dbQuery["name"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["serial_no"]; ok {
		dbQuery["serial_no"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["url"]; ok {
		dbQuery["url"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["description"]; ok {
		dbQuery["description"] = val[0]
	}

	if val, ok := filter["description_reg"]; ok {
		dbQuery["description"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["signal_title"]; ok {
		dbQuery["signal_title"] = val[0]
	}

	if val, ok := filter["signal_title_reg"]; ok {
		dbQuery["signal_title"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["signal_description"]; ok {
		dbQuery["signal_title"] = val[0]
	}

	if val, ok := filter["signal_description_reg"]; ok {
		dbQuery["signal_description"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	// const filters
	if val, ok := filter["hardware_type"]; ok {
		dbQuery["hardware_type"] = val[0]
	}

	if val, ok := filter["hardware_types"]; ok {
		dbQuery["hardware_type"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["hardware_model"]; ok {
		dbQuery["hardware_model"] = val[0]
	}

	if val, ok := filter["hardware_model"]; ok {
		dbQuery["hardware_models"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["hardware_vendor"]; ok {
		dbQuery["hardware_vendor"] = val[0]
	}

	if val, ok := filter["hardware_vendors"]; ok {
		dbQuery["hardware_vendor"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["type"]; ok {
		dbQuery["type"] = val[0]
	}

	if val, ok := filter["types"]; ok {
		dbQuery["type"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["category"]; ok {
		dbQuery["category"] = val[0]
	}

	if val, ok := filter["categories"]; ok {
		dbQuery["category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["status"]; ok {
		dbQuery["status"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["nin_status"]; ok {
		dbQuery["status"] = bsonMGO.M{"$nin": val}
	}

	if val, ok := filter["trigger"]; ok {
		dbQuery["trigger"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["is_approved"]; ok {
		b, parseErr := strconv.ParseBool(val[0])
		if parseErr == nil {
			dbQuery["is_approved"] = b
		}
	}

	if val, ok := filter["is_show_dmc"]; ok {
		b, parseErr := strconv.ParseBool(val[0])
		if parseErr == nil {
			dbQuery["is_show_dmc"] = b
		}
	}

	// alarm panel last status signal type
	if val, ok := filter["last_status_type"]; ok {
		dbQuery["last_status_signal.type"] = val[0]
	}

	if val, ok := filter["source_type"]; ok {
		dbQuery[fmt.Sprintf("%s.%s", "sources", val[0])] = filter["source_id"]
	}

	// config filters
	// alarm panel configuration
	if val, ok := filter["config_account"]; ok {
		dbQuery["config.account"] = val[0]
	}

	if val, ok := filter["config_phone_numbers"]; ok {
		dbQuery["config.phone_numbers"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["config_subscriber"]; ok {
		dbQuery["config.subscriber"] = val[0]
	}

	if val, ok := filter["config_signal_format"]; ok {
		dbQuery["config.signal_format"] = val[0]
	}

	if val, ok := filter["config_communication_module"]; ok {
		dbQuery["config.communication_module"] = val[0]
	}

	if val, ok := filter["config_communication_module_serial_no"]; ok {
		dbQuery["config.communication_module_serial_no"] = val[0]
	}

	// access control configuration
	if val, ok := filter["config_access_control_token"]; ok {
		dbQuery["config.access_control_token"] = val[0]
	}

	// access control logs list
	if val, ok := filter["personal_name"]; ok {
		dbQuery["personal_name"] = val[0]
	}

	// vehicle tracking configuration
	if val, ok := filter["config_vehicle_id"]; ok {
		dbQuery["config.vehicle_id"] = val[0]
	}

	// camera config
	if val, ok := filter["config_vguard_device_id"]; ok {
		dbQuery["config.vguard_device_id"] = val[0]
	}

	if val, ok := filter["config_ivideon_server_id"]; ok {
		dbQuery["config.ivideon_server_id"] = val[0]
	}

	if _, ok := filter["config_empty_ivideon_server_id"]; ok {
		dbQuery["config.ivideon_server_id"] = bsonMGO.M{
			"$exists": false,
		}
	}

	if val, ok := filter["config_is_getting_ivideon_service"]; ok {
		b, parseErr := strconv.ParseBool(val[0])
		if parseErr == nil {
			dbQuery["config.is_getting_ivideon_service"] = b
		}
	}

	if val, ok := filter["ivideon_account_email"]; ok {
		dbQuery["ivideon_account.email"] = val[0]
	}

	if val, ok := filter["ivideon_account_user_id"]; ok {
		dbQuery["ivideon_account.user_id"] = val[0]
	}

	// task messages
	if val, ok := filter["message_type"]; ok {
		dbQuery["messages.type"] = val[0]
	}

	if val, ok := filter["message_device_id"]; ok {
		if _, isEmbedded := filter["embedded"]; isEmbedded {
			dbQuery["messages.device.token"] = val[0]
		} else {
			dbQuery["messages.device_id"] = val[0]
		}
	}

	if val, ok := filter["device_distributor_id"]; ok {
		dbQuery["device.distributor_id"] = val[0]
	}

	if val, ok := filter["device_reseller_id"]; ok {
		dbQuery["device.reseller_id"] = val[0]
	}

	if val, ok := filter["device_account_id"]; ok {
		dbQuery["device.account_id"] = val[0]
	}

	if val, ok := filter["message_qualifier"]; ok {
		dbQuery["messages.qualifier"] = val[0]
	}

	if val, ok := filter["message_signal_categories"]; ok {
		dbQuery["messages.signal_category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["message_task_categories"]; ok {
		dbQuery["messages.task_category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["message_ticket_categories"]; ok {
		dbQuery["messages.ticket_category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["message_protocol"]; ok {
		dbQuery["messages.protocol"] = val[0]
	}

	if val, ok := filter["message_source"]; ok {
		dbQuery["messages.source"] = val[0]
	}

	if val, ok := filter["message_from"]; ok {
		dbQuery["messages.from"] = val[0]
	}

	if val, ok := filter["message_trigger"]; ok {
		dbQuery["messages.trigger"] = val[0]
	}

	if val, ok := filter["message_publish_id"]; ok {
		dbQuery["messages.publish_id"] = val[0]
	}

	// messages
	if val, ok := filter["qualifier"]; ok {
		dbQuery["qualifier"] = val[0]
	}

	if val, ok := filter["signal_categories"]; ok {
		dbQuery["signal_category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["task_categories"]; ok {
		dbQuery["task_category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["ticket_categories"]; ok {
		dbQuery["ticket_category"] = bsonMGO.M{"$in": val}
	}

	if val, ok := filter["protocol"]; ok {
		dbQuery["protocol"] = val[0]
	}

	if val, ok := filter["source"]; ok {
		dbQuery["source"] = val[0]
	}

	if val, ok := filter["from"]; ok {
		dbQuery["from"] = val[0]
	}

	if val, ok := filter["premise_city"]; ok {
		dbQuery["premise.address.city"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	if val, ok := filter["premise_province"]; ok {
		dbQuery["premise.address.province"] = bsonMGO.M{
			"$regex":   fmt.Sprintf(".*%v.*", val[0]),
			"$options": "i",
		}
	}

	// keyword
	if val, ok := filter["keyword"]; ok {
		dbQuery["$or"] = []bsonMGO.M{
			bsonMGO.M{
				"name": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
			bsonMGO.M{
				"first_name": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
			bsonMGO.M{
				"last_name": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
			bsonMGO.M{
				"email": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
			bsonMGO.M{
				"phone_number": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
			bsonMGO.M{
				"address.city": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
			bsonMGO.M{
				"address.province": bsonMGO.M{
					"$regex":   fmt.Sprintf(".*%v.*", val[0]),
					"$options": "i",
				},
			},
		}
	}

	return dbQuery, nil
}
