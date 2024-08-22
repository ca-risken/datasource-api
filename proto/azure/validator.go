package azure

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate for ListAzureDataSourceRequest
func (l *ListAzureDataSourceRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Name, validation.Length(0, 64)),
	)
}

// Validate for ListAzureRequest
func (l *ListAzureRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
		validation.Field(&l.SubscriptionId, validation.Length(0, 128)),
	)
}

// Validate for GetAzureRequest
func (g *GetAzureRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.AzureId, validation.Required),
	)
}

// Validate for PutAzureRequest
func (p *PutAzureRequest) Validate() error {
	if p.Azure == nil {
		return errors.New("required Azure")
	}
	if err := validation.ValidateStruct(p,
		validation.Field(&p.ProjectId, validation.Required, validation.In(p.Azure.ProjectId)),
	); err != nil {
		return err
	}
	return p.Azure.Validate()
}

// Validate for DeleteAzureRequest
func (d *DeleteAzureRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.AzureId, validation.Required),
	)
}

// Validate for ListRelAzureDataSourceRequest
func (l *ListRelAzureDataSourceRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.ProjectId, validation.Required),
	)
}

// Validate for GetRelAzureDataSourceRequest
func (g *GetRelAzureDataSourceRequest) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.AzureId, validation.Required),
		validation.Field(&g.AzureDataSourceId, validation.Required),
	)
}

// Validate for AttachRelAzureDataSourceRequest
func (a *AttachRelAzureDataSourceRequest) Validate() error {
	if a.RelAzureDataSource == nil {
		return errors.New("required RelAzureDataSource")
	}
	if err := validation.ValidateStruct(a,
		validation.Field(&a.ProjectId, validation.Required, validation.In(a.RelAzureDataSource.ProjectId)),
	); err != nil {
		return err
	}
	return a.RelAzureDataSource.Validate()
}

// Validate for DetachRelAzureDataSourceRequest
func (d *DetachRelAzureDataSourceRequest) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ProjectId, validation.Required),
		validation.Field(&d.AzureId, validation.Required),
		validation.Field(&d.AzureDataSourceId, validation.Required),
	)
}

// Validate for InvokeScanRequest
func (i *InvokeScanAzureRequest) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.ProjectId, validation.Required),
		validation.Field(&i.AzureId, validation.Required),
		validation.Field(&i.AzureDataSourceId, validation.Required),
	)
}

/**
 * Entity
**/

// Validate for AzureForUpsert
func (g *AzureForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.Name, validation.Required, validation.Length(0, 64)),
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.SubscriptionId, validation.Required, validation.Length(0, 128)),
		validation.Field(&g.VerificationCode, validation.Required, validation.Length(8, 128)),
	)
}

// Validate for RelAzureDataSourceForUpsert
func (g *RelAzureDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(g,
		validation.Field(&g.AzureDataSourceId, validation.Required),
		validation.Field(&g.ProjectId, validation.Required),
		validation.Field(&g.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
	)
}
