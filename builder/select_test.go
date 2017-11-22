package builder_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestSelect(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Select("test")
		is.Equal("SELECT test", query.String())
	}
	{
		query := loukoum.SelectDistinct("test")
		is.Equal("SELECT DISTINCT test", query.String())
	}
	{
		query := loukoum.Select(loukoum.Column("test").As("foobar"))
		is.Equal("SELECT test AS foobar", query.String())
	}
	{
		query := loukoum.Select("test", "foobar")
		is.Equal("SELECT test, foobar", query.String())
	}
	{
		query := loukoum.Select("test", loukoum.Column("test2").As("foobar"))
		is.Equal("SELECT test, test2 AS foobar", query.String())
	}
	{
		query := loukoum.Select("a", "b", loukoum.Column("c").As("x"))
		is.Equal("SELECT a, b, c AS x", query.String())
	}
	{
		query := loukoum.Select("a", loukoum.Column("b"), loukoum.Column("c").As("x"))
		is.Equal("SELECT a, b, c AS x", query.String())
	}
}

func TestSelect_From(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Select("a", "b", "c").From("foobar")
		is.Equal("SELECT a, b, c FROM foobar", query.String())
	}
	{
		query := loukoum.Select("a").From(loukoum.Table("foobar").As("example"))
		is.Equal("SELECT a FROM foobar AS example", query.String())
	}
}

func TestSelect_Join(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2 ON test1.id = test2.fk_id")

		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2", "test1.id = test2.fk_id")

		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2", "test1.id = test2.fk_id", loukoum.InnerJoin)

		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test3", "test3.fkey = test1.id", loukoum.LeftJoin)

		is.Equal("SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", "test4.gid = test2.id", loukoum.RightJoin)

		is.Equal("SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test5").
			Join("test3", "ON test3.id = test5.fk_id", loukoum.InnerJoin)

		is.Equal("SELECT a, b, c FROM test5 INNER JOIN test3 ON test3.id = test5.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", "test4.gid = test2.id").Join("test3", "test4.uid = test3.id")

		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", loukoum.On("test4.gid", "test2.id")).
			Join("test3", loukoum.On("test4.uid", "test3.id"))

		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join(loukoum.Table("test4"), loukoum.On("test4.gid", "test2.id")).
			Join(loukoum.Table("test3"), loukoum.On("test4.uid", "test3.id"))

		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
}

func TestSelect_WhereOperatorOrder(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1))

		is.Equal(`SELECT id FROM table WHERE (id = 1)`, query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (slug = 'foo'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			And(loukoum.Condition("title").Equal("hello"))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("slug").Equal("foo")).
			Or(loukoum.Condition("title").Equal("hello"))

		is.Equal("SELECT id FROM table WHERE (((id = 1) OR (slug = 'foo')) OR (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			Or(loukoum.Condition("title").Equal("hello"))

		is.Equal(`SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (title = 'hello'))`, query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.Or(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
			).
			Or(
				loukoum.And(loukoum.Condition("slug").Equal("foo"), loukoum.Condition("active").Equal(true)),
			)

		is.Equal(fmt.Sprint("SELECT id FROM table WHERE (((id = 1) OR (title = 'hello')) OR ",
			"((slug = 'foo') AND (active = true)))"), query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.And(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
			)

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("hello"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("hello")).
			Where(loukoum.Condition("disable").Equal(false))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (title = 'hello')) AND (disable = false))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(
				loukoum.Condition("slug").Equal("foo").And(loukoum.Condition("active").Equal(true)),
			)

		is.Equal("SELECT id FROM table WHERE ((id = 1) OR ((slug = 'foo') AND (active = true)))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1).And(loukoum.Condition("slug").Equal("foo"))).
			Or(loukoum.Condition("active").Equal(true))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (active = true))", query.String())
	}
}

func TestSelect_WhereEqual(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Equal(false))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (disabled = :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(false, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (disabled = false)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").NotEqual(false))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (disabled != :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(false, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (disabled != false)", query.String())
	}
}

func TestSelect_WhereIs(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Is(nil))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (disabled IS NULL)", stmt)
		is.Len(args, 0)

		is.Equal("SELECT id FROM table WHERE (disabled IS NULL)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("active").IsNot(true))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (active IS NOT :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(true, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (active IS NOT true)", query.String())
	}
}

func TestSelect_WhereGreaterThan(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThan(2))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (count > :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(2, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (count > 2)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThanOrEqual(4))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (count >= :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(4, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (count >= 4)", query.String())
	}
}

func TestSelect_WhereLessThan(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThan(3))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (count < :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(3, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (count < 3)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThanOrEqual(6))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (count <= :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(6, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (count <= 6)", query.String())
	}
}

func TestSelect_WhereLike(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").Like("foo%"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (title LIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (title LIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotLike("foo%"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (title NOT LIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (title NOT LIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").ILike("foo%"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (title ILIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (title ILIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotILike("foo%"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (title NOT ILIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (title NOT ILIKE 'foo%')", query.String())
	}
}

func TestSelect_WhereBetween(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").Between(10, 20))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (count BETWEEN :arg_1 AND :arg_2)", stmt)
		is.Len(args, 2)
		is.Equal(10, args[":arg_1"])
		is.Equal(20, args[":arg_2"])

		is.Equal("SELECT id FROM table WHERE (count BETWEEN 10 AND 20)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").NotBetween(50, 70))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (count NOT BETWEEN :arg_1 AND :arg_2)", stmt)
		is.Len(args, 2)
		is.Equal(50, args[":arg_1"])
		is.Equal(70, args[":arg_2"])

		is.Equal("SELECT id FROM table WHERE (count NOT BETWEEN 50 AND 70)", query.String())
	}
}

func TestSelect_WhereIn(t *testing.T) {
	is := require.New(t)

	// Slice of integers
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In([]int64{1, 2, 3}))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int64(1), args[":arg_1"])
		is.Equal(int64(2), args[":arg_2"])
		is.Equal(int64(3), args[":arg_3"])

		is.Equal("SELECT id FROM table WHERE (id IN (1, 2, 3))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn([]int{1, 2, 3}))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int(1), args[":arg_1"])
		is.Equal(int(2), args[":arg_2"])
		is.Equal(int(3), args[":arg_3"])

		is.Equal("SELECT id FROM table WHERE (id NOT IN (1, 2, 3))", query.String())
	}

	// Integers as variadic
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(1, 2, 3))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int(1), args[":arg_1"])
		is.Equal(int(2), args[":arg_2"])
		is.Equal(int(3), args[":arg_3"])

		is.Equal("SELECT id FROM table WHERE (id IN (1, 2, 3))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(1, 2, 3))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int(1), args[":arg_1"])
		is.Equal(int(2), args[":arg_2"])
		is.Equal(int(3), args[":arg_3"])

		is.Equal("SELECT id FROM table WHERE (id NOT IN (1, 2, 3))", query.String())
	}

	// Slice of strings
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"read", "unread"}))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args[":arg_1"])
		is.Equal("unread", args[":arg_2"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"read", "unread"}))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args[":arg_1"])
		is.Equal("unread", args[":arg_2"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))", query.String())
	}

	// Strings as variadic
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("read", "unread"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args[":arg_1"])
		is.Equal("unread", args[":arg_2"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("read", "unread"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args[":arg_1"])
		is.Equal("unread", args[":arg_2"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("read"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("read"))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"read"}))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"read"}))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read'))", query.String())
	}

	// Subquery
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Statement(),
			))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = :arg_1)))", stmt)
		is.Len(args, 1)
		is.Equal(1, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = 1)))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Statement(),
			))

		stmt, args := query.Prepare()
		is.Equal("SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = :arg_1)))", stmt)
		is.Len(args, 1)
		is.Equal(1, args[":arg_1"])

		is.Equal("SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = 1)))", query.String())
	}
}

func TestSelect_GroupBy(t *testing.T) {
	is := require.New(t)

	// One column
	{
		query := loukoum.
			Select("COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(true)).
			GroupBy("name")

		is.Equal("SELECT COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name", query.String())
	}
	{
		query := loukoum.
			Select("COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(true)).
			GroupBy(loukoum.Column("name"))

		is.Equal("SELECT COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name", query.String())
	}

	// Many columns
	{
		query := loukoum.
			Select("COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(true)).
			GroupBy("name", "email")

		is.Equal("SELECT COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name, email", query.String())
	}
	{
		query := loukoum.
			Select("COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(true)).
			GroupBy(loukoum.Column("name"), loukoum.Column("email"))

		is.Equal("SELECT COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name, email", query.String())
	}
	{
		query := loukoum.
			Select("COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(true)).
			GroupBy("name", "email", "user_id")

		is.Equal("SELECT COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name, email, user_id", query.String())
	}
	{
		query := loukoum.
			Select("COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(true)).
			GroupBy(loukoum.Column("name"), loukoum.Column("email"), loukoum.Column("user_id"))

		is.Equal("SELECT COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name, email, user_id", query.String())
	}
}