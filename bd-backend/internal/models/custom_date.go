package models

import (
    "database/sql/driver"
    "fmt"
    "time"
)

// CustomDate is a custom type for date format YYYY-MM-DD
type CustomDate time.Time

// Scan implements the sql.Scanner interface
func (cd *CustomDate) Scan(value interface{}) error {
    switch v := value.(type) {
    case time.Time:
        *cd = CustomDate(v)
        return nil
    case []byte:
        t, err := time.Parse("2006-01-02", string(v))
        if err != nil {
            return err
        }
        *cd = CustomDate(t)
        return nil
    case string:
        t, err := time.Parse("2006-01-02", v)
        if err != nil {
            return err
        }
        *cd = CustomDate(t)
        return nil
    default:
        return fmt.Errorf("cannot scan type %T into CustomDate", v)
    }
}

// Value implements the driver.Valuer interface
func (cd CustomDate) Value() (driver.Value, error) {
    return time.Time(cd).Format("2006-01-02"), nil
}
