-- TODO DELETE ME
SELECT 
	id,
	name,
	amount,
	date,
	icon,
	color
FROM (
	SELECT
		e.id,
		e.name,
		e.amount,
		e.date,
		e.user_id,
		c.icon,
		c.color
	FROM
		expense e
	JOIN category c ON c.id = e.category_id
	WHERE e.expense_group_id IS NULL
	UNION
	SELECT
		eg.id,
		eg.name,
		SUM(e.amount) AS amount,
		e.date,
		eg.user_id,
		c.icon,
		c.color
	FROM
		expense_group eg
	JOIN expense e ON e.expense_group_id = eg.id
	JOIN category c ON c.id = e.category_id
	GROUP BY eg.id
)
WHERE 
	date BETWEEN '2025-03-05' AND '2025-03-10'
	AND user_id = '1'
ORDER BY 
	date DESC,
	amount DESC
LIMIT 10
OFFSET 0;
