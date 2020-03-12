package v1

import (
	"context"
	"database/sql"
	"fmt"

	v1 "github.com/mark/todo/services/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
	deleteToDoQuery="DELETE FROM todolists WHERE `id`=?"
	selectOneQuery="SELECT * FROM todolists WHERE `id`=?"
	insertToDoQuery="INSERT INTO todolists(`title`, `description`, `order`) VALUES(?, ?, ?)"
	updateToDoQuery="UPDATE todolists SET `title`=?, `description`=?, `order`=? WHERE `id`=?"
)

// toDoServiceServer struct
type toDoServiceServer struct {
	db *sql.DB
}

// NewToDoServiceServer create new ToDoServiceServer instance
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{
		db: db,
	}
}

// checkAPI check if provided api version is valid
func (s *toDoServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'",
				apiVersion, api)
		}
	}
	return nil
}

// connect database context 
func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if  err != nil {
		return nil, status.Error(codes.Unknown, "Failed to connect to database -> "+err.Error())
	}
	return c, nil
}

// Create - create new `ToDo`
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	res, err := c.ExecContext(ctx, insertToDoQuery, req.ToDo.Title, req.ToDo.Description, req.ToDo.Order)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to insert into ToDo -> " + err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieved id for created todo -> " + err.Error())
	}
	return &v1.CreateResponse{
		Api: apiVersion,
		Id: id,
	}, nil
}

// Read - find ToDo object with provided `id`
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rows, err := c.QueryContext(ctx, selectOneQuery, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to select from Todo -> " + err.Error())
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "Failed to retrieve data from ToDo -> " + err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}

	var td v1.ToDo
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &td.Order); err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve field value from ToDo row -> " + err.Error())
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Found multiple ToDo rows with Id='%d'", req.Id))
	}
	return &v1.ReadResponse{
		Api: apiVersion,
		ToDo: &td,
	}, nil
}

// Update - update ToDo with provided `id`
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	res, err := c.ExecContext(ctx, updateToDoQuery, req.ToDo.Title, req.ToDo.Description, req.ToDo.Order, req.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to update ToDo -> " + err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve rows affected -> " + err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with Id='%d' is not found", req.ToDo.Id))
	}

	return &v1.UpdateResponse{
		Api: apiVersion,
		Updated: rows,
	}, nil
}

// Delete - delete ToDo with provided `id`
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		 return nil, err
	}
	defer c.Close()
	res, err := c.ExecContext(ctx, deleteToDoQuery, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to delete ToDo -> " + err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve rows affected value -> " + err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with Id='%d' is not found", req.Id))
	}
	return &v1.DeleteResponse{
		Api: apiVersion,
		Deleted: rows,
	}, nil
}

// ReadAll - read all todos list from ToDo
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	rows, err := c.QueryContext(ctx, "SELECT * FROM todolists")
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to select from ToDo->" + err.Error())
	}
	defer rows.Close()
	list := []*v1.ToDo{}
	for rows.Next() {
		td := new(v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &td.Order); err != nil {
			return nil, status.Error(codes.Unknown, "Failed to retrieve field value from ToDo row -> " + err.Error())
		}
		list = append(list, td)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve data from ToDo->" + err.Error())
	}
	return &v1.ReadAllResponse{
		Api: apiVersion,
		ToDos: list,
	}, nil
}