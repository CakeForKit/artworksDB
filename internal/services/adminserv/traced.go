package adminserv

import (
	"context"

	"github.com/CakeForKit/artworksDB.git/internal/models"
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedAdminService struct {
	adminService AdminService
	tracer       *tracing.Tracer
}

func NewTracedAdminService(adminService AdminService, tracer *tracing.Tracer) AdminService {
	return &tracedAdminService{
		adminService: adminService,
		tracer:       tracer,
	}
}

func (e *tracedAdminService) GetAllEmployees(ctx context.Context) ([]*models.Employee, error) {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "GetAllEmployees")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "adminserv"),
			attribute.String("role", "admin"),
		)
	}

	return e.adminService.GetAllEmployees(ctx)
}

func (e *tracedAdminService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "GetAllUsers")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "adminserv"),
			attribute.String("role", "admin"),
		)
	}

	return e.adminService.GetAllUsers(ctx)
}

func (e *tracedAdminService) ChangeEmployeeRights(ctx context.Context, employeeID uuid.UUID, valid bool) error {
	var span trace.Span
	if e.tracer.IsEnabled() {
		ctx, span = e.tracer.StartSpan(ctx, "ChangeEmployeeRights")
		defer span.End()

		span.SetAttributes(
			attribute.String("component", "adminserv"),
			attribute.String("employee_id", employeeID.String()),
			attribute.Bool("valid", valid),
			attribute.String("role", "admin"),
		)
	}

	return e.adminService.ChangeEmployeeRights(ctx, employeeID, valid)
}
