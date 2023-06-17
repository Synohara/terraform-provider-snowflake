package sdk

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestResourceMonitorCreate(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &CreateResourceMonitorOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE RESOURCE MONITOR WITH"
		assert.Equal(t, expected, actual)
	})

	t.Run("with complete options", func(t *testing.T) {
		creditQuota := Int(100)
		frequency := Monthly
		startTimeStamp := "IMMIEDIATELY"
		endTimeStamp := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).String()
		notifyUsers := []string{"FIRST_USER", "SECOND_USER"}
		triggers := []TriggerDefinition{
			{
				Threshold:     50,
				TriggerAction: SuspendImmediate,
			},
			{
				Threshold:     100,
				TriggerAction: Notify,
			},
		}

		opts := &CreateResourceMonitorOptions{
			OrReplace:      Bool(true),
			name:           id,
			CreditQuota:    creditQuota,
			Frequency:      &frequency,
			StartTimeStamp: &startTimeStamp,
			EndTimeStamp:   &endTimeStamp,
			// NotifyUsers:    &NotifyUsers{notifyUsers},
			Triggers:       &triggers,
		}

		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`CREATE OR REPLACE RESOURCE MONITOR %s WITH CREDIT_QUOTA = %d FREQUENCY = %s START_TIMESTAMP = %s END_TIMESTAMP = %s NOTIFY_USERS = (%s) TRIGGERS ON %d PERCENT DO %s ON %d PERCENT DO %s`,
			id.FullyQualifiedName(),
			*creditQuota,
			frequency,
			startTimeStamp,
			endTimeStamp,
			strings.Join(notifyUsers, " "),
			triggers[0].Threshold,
			triggers[0].TriggerAction,
			triggers[1].Threshold,
			triggers[1].TriggerAction,
		)

		assert.Equal(t, expected, actual)
	})
}

func TestResourceMonitorAlter(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "ALTER RESOURCE MONITOR"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &AlterResourceMonitorOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER RESOURCE MONITOR %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("with a single set", func(t *testing.T) {
		newCreditQuota := Int(50)
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				CreditQuota: newCreditQuota,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = %d", id.FullyQualifiedName(), *newCreditQuota)
		assert.Equal(t, expected, actual)
	})

	t.Run("with a mulitple set", func(t *testing.T) {
		newCreditQuota := Int(50)
		newFrequency := Yearly
		newStartTimeStamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).String()
		opts := &AlterResourceMonitorOptions{
			name: id,
			Set: &ResourceMonitorSet{
				CreditQuota:    newCreditQuota,
				Frequency:      &newFrequency,
				StartTimeStamp: &newStartTimeStamp,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER RESOURCE MONITOR %s SET CREDIT_QUOTA = %d FREQUENCY = %s START_TIMESTAMP = %s", id.FullyQualifiedName(), *newCreditQuota, newFrequency, newStartTimeStamp)
		assert.Equal(t, expected, actual)
	})
}

func TestResourceMonitorDrop(t *testing.T) {
	id := randomAccountObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &dropResourceMonitorOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "DROP RESOURCE MONITOR"
		assert.Equal(t, expected, actual)
	})

	t.Run("only name", func(t *testing.T) {
		opts := &dropResourceMonitorOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("DROP RESOURCE MONITOR %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestResourceMonitorShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &ShowResourceMonitorOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "SHOW RESOURCE MONITORS"
		assert.Equal(t, expected, actual)
	})

	t.Run("with like", func(t *testing.T) {
		opts := &ShowResourceMonitorOptions{
			Like: &Like{
				Pattern: String(id.Name()),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("SHOW RESOURCE MONITORS LIKE '%s'", id.Name())
		assert.Equal(t, expected, actual)
	})
}