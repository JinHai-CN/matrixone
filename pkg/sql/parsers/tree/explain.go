// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tree

import "strconv"

type Explain interface {
	Statement
}

type explainImpl struct {
	Explain
	Statement Statement
	Format    string
}

//EXPLAIN stmt statement
type ExplainStmt struct {
	explainImpl
}

func (node *ExplainStmt) Format(ctx *FmtCtx) {
	ctx.WriteString("explain")
	format := node.explainImpl.Format
	if format != "" && format != "row" {
		ctx.WriteString(" format = ")
		ctx.WriteString(node.explainImpl.Format)
	}
	stmt := node.explainImpl.Statement
	switch stmt.(type) {
	case *ShowColumns:
		st := stmt.(*ShowColumns)
		if st.Table != nil {
			ctx.WriteByte(' ')
			st.Table.Format(ctx)
		}
		if st.ColName != nil {
			ctx.WriteByte(' ')
			st.ColName.Format(ctx)
		}
	default:
		if stmt != nil {
			ctx.WriteByte(' ')
			stmt.Format(ctx)
		}
	}
}

func NewExplainStmt(stmt Statement, f string) *ExplainStmt {
	return &ExplainStmt{explainImpl{Statement: stmt, Format: f}}
}

//EXPLAIN ANALYZE statement
type ExplainAnalyze struct {
	explainImpl
}

func (node *ExplainAnalyze) Format(ctx *FmtCtx) {
	ctx.WriteString("explain analyze ")
	node.explainImpl.Statement.Format(ctx)
}

func NewExplainAnalyze(stmt Statement, f string) *ExplainAnalyze {
	return &ExplainAnalyze{explainImpl{Statement: stmt, Format: f}}
}

//EXPLAIN FOR CONNECTION statement
type ExplainFor struct {
	explainImpl
	ID uint64
}

func (node *ExplainFor) Format(ctx *FmtCtx) {
	ctx.WriteString("explain format = ")
	ctx.WriteString(node.explainImpl.Format)
	ctx.WriteString(" for connection ")
	ctx.WriteString(strconv.FormatInt(int64(node.ID), 10))
}

func NewExplainFor(f string, id uint64) *ExplainFor {
	return &ExplainFor{
		explainImpl: explainImpl{Statement: nil, Format: f},
		ID:          id,
	}
}
