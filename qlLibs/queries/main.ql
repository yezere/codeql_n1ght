/**
 * @name N1ght QL Inspector
 * @kind problem
 * @description 查找所有Java方法的测试查询
 * @problem.severity warning
 * @precision high
 * @id java/n1ght-ql-inspector
 * @tags security
 *       maintainability
 */
import java

from Method m
where m.getName() = "exec"
select m, "N1ght QL Inspector: Found method " + m.getName()