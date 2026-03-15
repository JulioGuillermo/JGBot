package taskinfrastructure

import (
	"testing"
	"time"
)

func TestCronMngNew(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	if cronMng == nil {
		t.Fatal("NewCronMng() returned nil")
	}

	if cronMng.Cron == nil {
		t.Fatal("NewCronMng() Cron field is nil")
	}
}

func TestCronMngAddFunc(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	cmd := func() {}

	id, err := cronMng.AddFunc("* * * * *", cmd)
	if err != nil {
		t.Fatalf("AddFunc() error = %v", err)
	}

	if id <= 0 {
		t.Errorf("AddFunc() returned invalid id = %d, want > 0", id)
	}

	// Give some time for the cron to potentially execute
	time.Sleep(100 * time.Millisecond)
}

func TestCronMngRemove(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	cmd := func() {}

	id, err := cronMng.AddFunc("* * * * *", cmd)
	if err != nil {
		t.Fatalf("AddFunc() error = %v", err)
	}

	cronMng.Remove(id)

	// The entry should be removed, so the cmd should not execute
	// We can't directly verify this, but we can ensure no panic occurs
	time.Sleep(100 * time.Millisecond)
}

func TestCronMngAddFuncInvalidSpec(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	cmd := func() {}

	_, err := cronMng.AddFunc("invalid spec", cmd)
	if err == nil {
		t.Error("AddFunc() with invalid spec should return error")
	}
}

func TestCronMngMultipleEntries(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	cmd := func() {}

	// Add multiple entries
	ids := make([]int, 5)
	for i := range 5 {
		id, err := cronMng.AddFunc("* * * * *", cmd)
		if err != nil {
			t.Fatalf("AddFunc() error = %v", err)
		}
		ids[i] = id
	}

	// Verify all IDs are unique
	seen := make(map[int]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		seen[id] = true
	}

	// Remove some entries
	cronMng.Remove(ids[0])
	cronMng.Remove(ids[2])
	cronMng.Remove(ids[4])

	time.Sleep(100 * time.Millisecond)
}

func TestCronMngRemoveNonExistent(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	// Removing a non-existent entry should not panic
	cronMng.Remove(9999)
}

func TestCronMngAddFuncWithDifferentSchedules(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	tests := []struct {
		name    string
		spec    string
		wantErr bool
	}{
		{
			name:    "every minute",
			spec:    "* * * * *",
			wantErr: false,
		},
		{
			name:    "every hour",
			spec:    "0 * * * *",
			wantErr: false,
		},
		{
			name:    "every day at midnight",
			spec:    "0 0 * * *",
			wantErr: false,
		},
		{
			name:    "every monday at 9am",
			spec:    "0 9 * * 1",
			wantErr: false,
		},
		{
			name:    "specific date and time",
			spec:    "0 0 1 1 *",
			wantErr: false,
		},
		{
			name:    "invalid spec",
			spec:    "invalid",
			wantErr: true,
		},
		{
			name:    "empty spec",
			spec:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd := func() {}
			_, err := cronMng.AddFunc(tt.spec, cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCronMngInterfaceCompliance(t *testing.T) {
	t.Parallel()

	// Verify that CronMng implements taskdomain.CronMng interface
	var _ interface {
		AddFunc(string, func()) (int, error)
		Remove(int)
	} = &CronMng{}
}

func TestCronMngUnderlyingCron(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	// Verify the underlying cron is properly initialized
	if cronMng.Cron == nil {
		t.Fatal("Underlying Cron should not be nil")
	}

	// Verify we can access cron entries
	entries := cronMng.Entries()
	if entries == nil {
		t.Error("Entries() should return a slice, not nil")
	}
}

func TestCronMngConcurrentAccess(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	done := make(chan bool, 10)

	// Concurrent AddFunc calls
	for range 5 {
		go func() {
			defer func() { done <- true }()
			cmd := func() {}
			_, err := cronMng.AddFunc("* * * * *", cmd)
			if err != nil {
				t.Errorf("Concurrent AddFunc() error = %v", err)
			}
		}()
	}

	// Wait for some AddFunc calls to complete
	time.Sleep(50 * time.Millisecond)

	// Concurrent Remove calls
	for i := range 5 {
		go func(id int) {
			defer func() { done <- true }()
			cronMng.Remove(id)
		}(i + 1)
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}
}

func TestCronMngEntryIDConversion(t *testing.T) {
	t.Parallel()

	cronMng := NewCronMng()

	cmd := func() {}

	id, err := cronMng.AddFunc("* * * * *", cmd)
	if err != nil {
		t.Fatalf("AddFunc() error = %v", err)
	}

	// Verify the ID can be used with Remove
	cronMng.Remove(id)

	// Should not panic
	time.Sleep(50 * time.Millisecond)
}

func TestCronMngStress(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	cronMng := NewCronMng()

	const numEntries = 100
	ids := make([]int, numEntries)

	// Add many entries
	for i := range numEntries {
		cmd := func() {}
		id, err := cronMng.AddFunc("* * * * *", cmd)
		if err != nil {
			t.Fatalf("AddFunc() error = %v", err)
		}
		ids[i] = id
	}

	// Remove all entries
	for _, id := range ids {
		cronMng.Remove(id)
	}

	time.Sleep(100 * time.Millisecond)
}
