// Code generated from /var/antlrResult/Modelica.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // Modelica

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseModelicaListener is a complete listener for a parse tree produced by ModelicaParser.
type BaseModelicaListener struct{}

var _ ModelicaListener = &BaseModelicaListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseModelicaListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseModelicaListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseModelicaListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseModelicaListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStored_definition is called when production stored_definition is entered.
func (s *BaseModelicaListener) EnterStored_definition(ctx *Stored_definitionContext) {}

// ExitStored_definition is called when production stored_definition is exited.
func (s *BaseModelicaListener) ExitStored_definition(ctx *Stored_definitionContext) {}

// EnterClass_definition is called when production class_definition is entered.
func (s *BaseModelicaListener) EnterClass_definition(ctx *Class_definitionContext) {}

// ExitClass_definition is called when production class_definition is exited.
func (s *BaseModelicaListener) ExitClass_definition(ctx *Class_definitionContext) {}

// EnterClass_specifier is called when production class_specifier is entered.
func (s *BaseModelicaListener) EnterClass_specifier(ctx *Class_specifierContext) {}

// ExitClass_specifier is called when production class_specifier is exited.
func (s *BaseModelicaListener) ExitClass_specifier(ctx *Class_specifierContext) {}

// EnterClass_prefixes is called when production class_prefixes is entered.
func (s *BaseModelicaListener) EnterClass_prefixes(ctx *Class_prefixesContext) {}

// ExitClass_prefixes is called when production class_prefixes is exited.
func (s *BaseModelicaListener) ExitClass_prefixes(ctx *Class_prefixesContext) {}

// EnterLong_class_specifier is called when production long_class_specifier is entered.
func (s *BaseModelicaListener) EnterLong_class_specifier(ctx *Long_class_specifierContext) {}

// ExitLong_class_specifier is called when production long_class_specifier is exited.
func (s *BaseModelicaListener) ExitLong_class_specifier(ctx *Long_class_specifierContext) {}

// EnterShort_class_specifier is called when production short_class_specifier is entered.
func (s *BaseModelicaListener) EnterShort_class_specifier(ctx *Short_class_specifierContext) {}

// ExitShort_class_specifier is called when production short_class_specifier is exited.
func (s *BaseModelicaListener) ExitShort_class_specifier(ctx *Short_class_specifierContext) {}

// EnterDer_class_specifier is called when production der_class_specifier is entered.
func (s *BaseModelicaListener) EnterDer_class_specifier(ctx *Der_class_specifierContext) {}

// ExitDer_class_specifier is called when production der_class_specifier is exited.
func (s *BaseModelicaListener) ExitDer_class_specifier(ctx *Der_class_specifierContext) {}

// EnterBase_prefix is called when production base_prefix is entered.
func (s *BaseModelicaListener) EnterBase_prefix(ctx *Base_prefixContext) {}

// ExitBase_prefix is called when production base_prefix is exited.
func (s *BaseModelicaListener) ExitBase_prefix(ctx *Base_prefixContext) {}

// EnterEnum_list is called when production enum_list is entered.
func (s *BaseModelicaListener) EnterEnum_list(ctx *Enum_listContext) {}

// ExitEnum_list is called when production enum_list is exited.
func (s *BaseModelicaListener) ExitEnum_list(ctx *Enum_listContext) {}

// EnterEnumeration_literal is called when production enumeration_literal is entered.
func (s *BaseModelicaListener) EnterEnumeration_literal(ctx *Enumeration_literalContext) {}

// ExitEnumeration_literal is called when production enumeration_literal is exited.
func (s *BaseModelicaListener) ExitEnumeration_literal(ctx *Enumeration_literalContext) {}

// EnterComposition is called when production composition is entered.
func (s *BaseModelicaListener) EnterComposition(ctx *CompositionContext) {}

// ExitComposition is called when production composition is exited.
func (s *BaseModelicaListener) ExitComposition(ctx *CompositionContext) {}

// EnterModel_annotation is called when production model_annotation is entered.
func (s *BaseModelicaListener) EnterModel_annotation(ctx *Model_annotationContext) {}

// ExitModel_annotation is called when production model_annotation is exited.
func (s *BaseModelicaListener) ExitModel_annotation(ctx *Model_annotationContext) {}

// EnterLanguage_specification is called when production language_specification is entered.
func (s *BaseModelicaListener) EnterLanguage_specification(ctx *Language_specificationContext) {}

// ExitLanguage_specification is called when production language_specification is exited.
func (s *BaseModelicaListener) ExitLanguage_specification(ctx *Language_specificationContext) {}

// EnterExternal_function_call is called when production external_function_call is entered.
func (s *BaseModelicaListener) EnterExternal_function_call(ctx *External_function_callContext) {}

// ExitExternal_function_call is called when production external_function_call is exited.
func (s *BaseModelicaListener) ExitExternal_function_call(ctx *External_function_callContext) {}

// EnterElement_list is called when production element_list is entered.
func (s *BaseModelicaListener) EnterElement_list(ctx *Element_listContext) {}

// ExitElement_list is called when production element_list is exited.
func (s *BaseModelicaListener) ExitElement_list(ctx *Element_listContext) {}

// EnterElement is called when production element is entered.
func (s *BaseModelicaListener) EnterElement(ctx *ElementContext) {}

// ExitElement is called when production element is exited.
func (s *BaseModelicaListener) ExitElement(ctx *ElementContext) {}

// EnterImport_clause is called when production import_clause is entered.
func (s *BaseModelicaListener) EnterImport_clause(ctx *Import_clauseContext) {}

// ExitImport_clause is called when production import_clause is exited.
func (s *BaseModelicaListener) ExitImport_clause(ctx *Import_clauseContext) {}

// EnterImport_list is called when production import_list is entered.
func (s *BaseModelicaListener) EnterImport_list(ctx *Import_listContext) {}

// ExitImport_list is called when production import_list is exited.
func (s *BaseModelicaListener) ExitImport_list(ctx *Import_listContext) {}

// EnterExtends_clause is called when production extends_clause is entered.
func (s *BaseModelicaListener) EnterExtends_clause(ctx *Extends_clauseContext) {}

// ExitExtends_clause is called when production extends_clause is exited.
func (s *BaseModelicaListener) ExitExtends_clause(ctx *Extends_clauseContext) {}

// EnterConstraining_clause is called when production constraining_clause is entered.
func (s *BaseModelicaListener) EnterConstraining_clause(ctx *Constraining_clauseContext) {}

// ExitConstraining_clause is called when production constraining_clause is exited.
func (s *BaseModelicaListener) ExitConstraining_clause(ctx *Constraining_clauseContext) {}

// EnterComponent_clause is called when production component_clause is entered.
func (s *BaseModelicaListener) EnterComponent_clause(ctx *Component_clauseContext) {}

// ExitComponent_clause is called when production component_clause is exited.
func (s *BaseModelicaListener) ExitComponent_clause(ctx *Component_clauseContext) {}

// EnterType_prefix is called when production type_prefix is entered.
func (s *BaseModelicaListener) EnterType_prefix(ctx *Type_prefixContext) {}

// ExitType_prefix is called when production type_prefix is exited.
func (s *BaseModelicaListener) ExitType_prefix(ctx *Type_prefixContext) {}

// EnterType_specifier is called when production type_specifier is entered.
func (s *BaseModelicaListener) EnterType_specifier(ctx *Type_specifierContext) {}

// ExitType_specifier is called when production type_specifier is exited.
func (s *BaseModelicaListener) ExitType_specifier(ctx *Type_specifierContext) {}

// EnterComponent_list is called when production component_list is entered.
func (s *BaseModelicaListener) EnterComponent_list(ctx *Component_listContext) {}

// ExitComponent_list is called when production component_list is exited.
func (s *BaseModelicaListener) ExitComponent_list(ctx *Component_listContext) {}

// EnterComponent_declaration is called when production component_declaration is entered.
func (s *BaseModelicaListener) EnterComponent_declaration(ctx *Component_declarationContext) {}

// ExitComponent_declaration is called when production component_declaration is exited.
func (s *BaseModelicaListener) ExitComponent_declaration(ctx *Component_declarationContext) {}

// EnterCondition_attribute is called when production condition_attribute is entered.
func (s *BaseModelicaListener) EnterCondition_attribute(ctx *Condition_attributeContext) {}

// ExitCondition_attribute is called when production condition_attribute is exited.
func (s *BaseModelicaListener) ExitCondition_attribute(ctx *Condition_attributeContext) {}

// EnterDeclaration is called when production declaration is entered.
func (s *BaseModelicaListener) EnterDeclaration(ctx *DeclarationContext) {}

// ExitDeclaration is called when production declaration is exited.
func (s *BaseModelicaListener) ExitDeclaration(ctx *DeclarationContext) {}

// EnterModification is called when production modification is entered.
func (s *BaseModelicaListener) EnterModification(ctx *ModificationContext) {}

// ExitModification is called when production modification is exited.
func (s *BaseModelicaListener) ExitModification(ctx *ModificationContext) {}

// EnterClass_modification is called when production class_modification is entered.
func (s *BaseModelicaListener) EnterClass_modification(ctx *Class_modificationContext) {}

// ExitClass_modification is called when production class_modification is exited.
func (s *BaseModelicaListener) ExitClass_modification(ctx *Class_modificationContext) {}

// EnterArgument_list is called when production argument_list is entered.
func (s *BaseModelicaListener) EnterArgument_list(ctx *Argument_listContext) {}

// ExitArgument_list is called when production argument_list is exited.
func (s *BaseModelicaListener) ExitArgument_list(ctx *Argument_listContext) {}

// EnterArgument is called when production argument is entered.
func (s *BaseModelicaListener) EnterArgument(ctx *ArgumentContext) {}

// ExitArgument is called when production argument is exited.
func (s *BaseModelicaListener) ExitArgument(ctx *ArgumentContext) {}

// EnterElement_modification_or_replaceable is called when production element_modification_or_replaceable is entered.
func (s *BaseModelicaListener) EnterElement_modification_or_replaceable(ctx *Element_modification_or_replaceableContext) {}

// ExitElement_modification_or_replaceable is called when production element_modification_or_replaceable is exited.
func (s *BaseModelicaListener) ExitElement_modification_or_replaceable(ctx *Element_modification_or_replaceableContext) {}

// EnterElement_modification is called when production element_modification is entered.
func (s *BaseModelicaListener) EnterElement_modification(ctx *Element_modificationContext) {}

// ExitElement_modification is called when production element_modification is exited.
func (s *BaseModelicaListener) ExitElement_modification(ctx *Element_modificationContext) {}

// EnterElement_redeclaration is called when production element_redeclaration is entered.
func (s *BaseModelicaListener) EnterElement_redeclaration(ctx *Element_redeclarationContext) {}

// ExitElement_redeclaration is called when production element_redeclaration is exited.
func (s *BaseModelicaListener) ExitElement_redeclaration(ctx *Element_redeclarationContext) {}

// EnterElement_replaceable is called when production element_replaceable is entered.
func (s *BaseModelicaListener) EnterElement_replaceable(ctx *Element_replaceableContext) {}

// ExitElement_replaceable is called when production element_replaceable is exited.
func (s *BaseModelicaListener) ExitElement_replaceable(ctx *Element_replaceableContext) {}

// EnterComponent_clause1 is called when production component_clause1 is entered.
func (s *BaseModelicaListener) EnterComponent_clause1(ctx *Component_clause1Context) {}

// ExitComponent_clause1 is called when production component_clause1 is exited.
func (s *BaseModelicaListener) ExitComponent_clause1(ctx *Component_clause1Context) {}

// EnterComponent_declaration1 is called when production component_declaration1 is entered.
func (s *BaseModelicaListener) EnterComponent_declaration1(ctx *Component_declaration1Context) {}

// ExitComponent_declaration1 is called when production component_declaration1 is exited.
func (s *BaseModelicaListener) ExitComponent_declaration1(ctx *Component_declaration1Context) {}

// EnterShort_class_definition is called when production short_class_definition is entered.
func (s *BaseModelicaListener) EnterShort_class_definition(ctx *Short_class_definitionContext) {}

// ExitShort_class_definition is called when production short_class_definition is exited.
func (s *BaseModelicaListener) ExitShort_class_definition(ctx *Short_class_definitionContext) {}

// EnterEquation_section is called when production equation_section is entered.
func (s *BaseModelicaListener) EnterEquation_section(ctx *Equation_sectionContext) {}

// ExitEquation_section is called when production equation_section is exited.
func (s *BaseModelicaListener) ExitEquation_section(ctx *Equation_sectionContext) {}

// EnterEquations is called when production equations is entered.
func (s *BaseModelicaListener) EnterEquations(ctx *EquationsContext) {}

// ExitEquations is called when production equations is exited.
func (s *BaseModelicaListener) ExitEquations(ctx *EquationsContext) {}

// EnterAlgorithm_section is called when production algorithm_section is entered.
func (s *BaseModelicaListener) EnterAlgorithm_section(ctx *Algorithm_sectionContext) {}

// ExitAlgorithm_section is called when production algorithm_section is exited.
func (s *BaseModelicaListener) ExitAlgorithm_section(ctx *Algorithm_sectionContext) {}

// EnterAlgorithm_statements is called when production algorithm_statements is entered.
func (s *BaseModelicaListener) EnterAlgorithm_statements(ctx *Algorithm_statementsContext) {}

// ExitAlgorithm_statements is called when production algorithm_statements is exited.
func (s *BaseModelicaListener) ExitAlgorithm_statements(ctx *Algorithm_statementsContext) {}

// EnterEquation is called when production equation is entered.
func (s *BaseModelicaListener) EnterEquation(ctx *EquationContext) {}

// ExitEquation is called when production equation is exited.
func (s *BaseModelicaListener) ExitEquation(ctx *EquationContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseModelicaListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseModelicaListener) ExitStatement(ctx *StatementContext) {}

// EnterIf_equation is called when production if_equation is entered.
func (s *BaseModelicaListener) EnterIf_equation(ctx *If_equationContext) {}

// ExitIf_equation is called when production if_equation is exited.
func (s *BaseModelicaListener) ExitIf_equation(ctx *If_equationContext) {}

// EnterIf_statement is called when production if_statement is entered.
func (s *BaseModelicaListener) EnterIf_statement(ctx *If_statementContext) {}

// ExitIf_statement is called when production if_statement is exited.
func (s *BaseModelicaListener) ExitIf_statement(ctx *If_statementContext) {}

// EnterControl_structure_body is called when production control_structure_body is entered.
func (s *BaseModelicaListener) EnterControl_structure_body(ctx *Control_structure_bodyContext) {}

// ExitControl_structure_body is called when production control_structure_body is exited.
func (s *BaseModelicaListener) ExitControl_structure_body(ctx *Control_structure_bodyContext) {}

// EnterFor_equation is called when production for_equation is entered.
func (s *BaseModelicaListener) EnterFor_equation(ctx *For_equationContext) {}

// ExitFor_equation is called when production for_equation is exited.
func (s *BaseModelicaListener) ExitFor_equation(ctx *For_equationContext) {}

// EnterFor_statement is called when production for_statement is entered.
func (s *BaseModelicaListener) EnterFor_statement(ctx *For_statementContext) {}

// ExitFor_statement is called when production for_statement is exited.
func (s *BaseModelicaListener) ExitFor_statement(ctx *For_statementContext) {}

// EnterFor_indices is called when production for_indices is entered.
func (s *BaseModelicaListener) EnterFor_indices(ctx *For_indicesContext) {}

// ExitFor_indices is called when production for_indices is exited.
func (s *BaseModelicaListener) ExitFor_indices(ctx *For_indicesContext) {}

// EnterFor_index is called when production for_index is entered.
func (s *BaseModelicaListener) EnterFor_index(ctx *For_indexContext) {}

// ExitFor_index is called when production for_index is exited.
func (s *BaseModelicaListener) ExitFor_index(ctx *For_indexContext) {}

// EnterWhile_statement is called when production while_statement is entered.
func (s *BaseModelicaListener) EnterWhile_statement(ctx *While_statementContext) {}

// ExitWhile_statement is called when production while_statement is exited.
func (s *BaseModelicaListener) ExitWhile_statement(ctx *While_statementContext) {}

// EnterWhen_equation is called when production when_equation is entered.
func (s *BaseModelicaListener) EnterWhen_equation(ctx *When_equationContext) {}

// ExitWhen_equation is called when production when_equation is exited.
func (s *BaseModelicaListener) ExitWhen_equation(ctx *When_equationContext) {}

// EnterWhen_statement is called when production when_statement is entered.
func (s *BaseModelicaListener) EnterWhen_statement(ctx *When_statementContext) {}

// ExitWhen_statement is called when production when_statement is exited.
func (s *BaseModelicaListener) ExitWhen_statement(ctx *When_statementContext) {}

// EnterConnect_clause is called when production connect_clause is entered.
func (s *BaseModelicaListener) EnterConnect_clause(ctx *Connect_clauseContext) {}

// ExitConnect_clause is called when production connect_clause is exited.
func (s *BaseModelicaListener) ExitConnect_clause(ctx *Connect_clauseContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseModelicaListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseModelicaListener) ExitExpression(ctx *ExpressionContext) {}

// EnterSimple_expression is called when production simple_expression is entered.
func (s *BaseModelicaListener) EnterSimple_expression(ctx *Simple_expressionContext) {}

// ExitSimple_expression is called when production simple_expression is exited.
func (s *BaseModelicaListener) ExitSimple_expression(ctx *Simple_expressionContext) {}

// EnterIf_expression is called when production if_expression is entered.
func (s *BaseModelicaListener) EnterIf_expression(ctx *If_expressionContext) {}

// ExitIf_expression is called when production if_expression is exited.
func (s *BaseModelicaListener) ExitIf_expression(ctx *If_expressionContext) {}

// EnterIf_expression_body is called when production if_expression_body is entered.
func (s *BaseModelicaListener) EnterIf_expression_body(ctx *If_expression_bodyContext) {}

// ExitIf_expression_body is called when production if_expression_body is exited.
func (s *BaseModelicaListener) ExitIf_expression_body(ctx *If_expression_bodyContext) {}

// EnterIf_expression_condition is called when production if_expression_condition is entered.
func (s *BaseModelicaListener) EnterIf_expression_condition(ctx *If_expression_conditionContext) {}

// ExitIf_expression_condition is called when production if_expression_condition is exited.
func (s *BaseModelicaListener) ExitIf_expression_condition(ctx *If_expression_conditionContext) {}

// EnterElseif_expression_condition is called when production elseif_expression_condition is entered.
func (s *BaseModelicaListener) EnterElseif_expression_condition(ctx *Elseif_expression_conditionContext) {}

// ExitElseif_expression_condition is called when production elseif_expression_condition is exited.
func (s *BaseModelicaListener) ExitElseif_expression_condition(ctx *Elseif_expression_conditionContext) {}

// EnterElse_expression_condition is called when production else_expression_condition is entered.
func (s *BaseModelicaListener) EnterElse_expression_condition(ctx *Else_expression_conditionContext) {}

// ExitElse_expression_condition is called when production else_expression_condition is exited.
func (s *BaseModelicaListener) ExitElse_expression_condition(ctx *Else_expression_conditionContext) {}

// EnterLogical_expression is called when production logical_expression is entered.
func (s *BaseModelicaListener) EnterLogical_expression(ctx *Logical_expressionContext) {}

// ExitLogical_expression is called when production logical_expression is exited.
func (s *BaseModelicaListener) ExitLogical_expression(ctx *Logical_expressionContext) {}

// EnterLogical_term is called when production logical_term is entered.
func (s *BaseModelicaListener) EnterLogical_term(ctx *Logical_termContext) {}

// ExitLogical_term is called when production logical_term is exited.
func (s *BaseModelicaListener) ExitLogical_term(ctx *Logical_termContext) {}

// EnterLogical_factor is called when production logical_factor is entered.
func (s *BaseModelicaListener) EnterLogical_factor(ctx *Logical_factorContext) {}

// ExitLogical_factor is called when production logical_factor is exited.
func (s *BaseModelicaListener) ExitLogical_factor(ctx *Logical_factorContext) {}

// EnterRelation is called when production relation is entered.
func (s *BaseModelicaListener) EnterRelation(ctx *RelationContext) {}

// ExitRelation is called when production relation is exited.
func (s *BaseModelicaListener) ExitRelation(ctx *RelationContext) {}

// EnterRel_op is called when production rel_op is entered.
func (s *BaseModelicaListener) EnterRel_op(ctx *Rel_opContext) {}

// ExitRel_op is called when production rel_op is exited.
func (s *BaseModelicaListener) ExitRel_op(ctx *Rel_opContext) {}

// EnterArithmetic_expression is called when production arithmetic_expression is entered.
func (s *BaseModelicaListener) EnterArithmetic_expression(ctx *Arithmetic_expressionContext) {}

// ExitArithmetic_expression is called when production arithmetic_expression is exited.
func (s *BaseModelicaListener) ExitArithmetic_expression(ctx *Arithmetic_expressionContext) {}

// EnterAdd_op is called when production add_op is entered.
func (s *BaseModelicaListener) EnterAdd_op(ctx *Add_opContext) {}

// ExitAdd_op is called when production add_op is exited.
func (s *BaseModelicaListener) ExitAdd_op(ctx *Add_opContext) {}

// EnterTerm is called when production term is entered.
func (s *BaseModelicaListener) EnterTerm(ctx *TermContext) {}

// ExitTerm is called when production term is exited.
func (s *BaseModelicaListener) ExitTerm(ctx *TermContext) {}

// EnterMul_op is called when production mul_op is entered.
func (s *BaseModelicaListener) EnterMul_op(ctx *Mul_opContext) {}

// ExitMul_op is called when production mul_op is exited.
func (s *BaseModelicaListener) ExitMul_op(ctx *Mul_opContext) {}

// EnterFactor is called when production factor is entered.
func (s *BaseModelicaListener) EnterFactor(ctx *FactorContext) {}

// ExitFactor is called when production factor is exited.
func (s *BaseModelicaListener) ExitFactor(ctx *FactorContext) {}

// EnterPrimary is called when production primary is entered.
func (s *BaseModelicaListener) EnterPrimary(ctx *PrimaryContext) {}

// ExitPrimary is called when production primary is exited.
func (s *BaseModelicaListener) ExitPrimary(ctx *PrimaryContext) {}

// EnterVector is called when production vector is entered.
func (s *BaseModelicaListener) EnterVector(ctx *VectorContext) {}

// ExitVector is called when production vector is exited.
func (s *BaseModelicaListener) ExitVector(ctx *VectorContext) {}

// EnterArray_arguments is called when production array_arguments is entered.
func (s *BaseModelicaListener) EnterArray_arguments(ctx *Array_argumentsContext) {}

// ExitArray_arguments is called when production array_arguments is exited.
func (s *BaseModelicaListener) ExitArray_arguments(ctx *Array_argumentsContext) {}

// EnterArray_iterator_constructor is called when production array_iterator_constructor is entered.
func (s *BaseModelicaListener) EnterArray_iterator_constructor(ctx *Array_iterator_constructorContext) {}

// ExitArray_iterator_constructor is called when production array_iterator_constructor is exited.
func (s *BaseModelicaListener) ExitArray_iterator_constructor(ctx *Array_iterator_constructorContext) {}

// EnterName is called when production name is entered.
func (s *BaseModelicaListener) EnterName(ctx *NameContext) {}

// ExitName is called when production name is exited.
func (s *BaseModelicaListener) ExitName(ctx *NameContext) {}

// EnterComponent_reference is called when production component_reference is entered.
func (s *BaseModelicaListener) EnterComponent_reference(ctx *Component_referenceContext) {}

// ExitComponent_reference is called when production component_reference is exited.
func (s *BaseModelicaListener) ExitComponent_reference(ctx *Component_referenceContext) {}

// EnterFunction_call_args is called when production function_call_args is entered.
func (s *BaseModelicaListener) EnterFunction_call_args(ctx *Function_call_argsContext) {}

// ExitFunction_call_args is called when production function_call_args is exited.
func (s *BaseModelicaListener) ExitFunction_call_args(ctx *Function_call_argsContext) {}

// EnterFunction_arguments is called when production function_arguments is entered.
func (s *BaseModelicaListener) EnterFunction_arguments(ctx *Function_argumentsContext) {}

// ExitFunction_arguments is called when production function_arguments is exited.
func (s *BaseModelicaListener) ExitFunction_arguments(ctx *Function_argumentsContext) {}

// EnterNamed_arguments is called when production named_arguments is entered.
func (s *BaseModelicaListener) EnterNamed_arguments(ctx *Named_argumentsContext) {}

// ExitNamed_arguments is called when production named_arguments is exited.
func (s *BaseModelicaListener) ExitNamed_arguments(ctx *Named_argumentsContext) {}

// EnterNamed_argument is called when production named_argument is entered.
func (s *BaseModelicaListener) EnterNamed_argument(ctx *Named_argumentContext) {}

// ExitNamed_argument is called when production named_argument is exited.
func (s *BaseModelicaListener) ExitNamed_argument(ctx *Named_argumentContext) {}

// EnterFunction_argument is called when production function_argument is entered.
func (s *BaseModelicaListener) EnterFunction_argument(ctx *Function_argumentContext) {}

// ExitFunction_argument is called when production function_argument is exited.
func (s *BaseModelicaListener) ExitFunction_argument(ctx *Function_argumentContext) {}

// EnterOutput_expression_list is called when production output_expression_list is entered.
func (s *BaseModelicaListener) EnterOutput_expression_list(ctx *Output_expression_listContext) {}

// ExitOutput_expression_list is called when production output_expression_list is exited.
func (s *BaseModelicaListener) ExitOutput_expression_list(ctx *Output_expression_listContext) {}

// EnterExpression_list is called when production expression_list is entered.
func (s *BaseModelicaListener) EnterExpression_list(ctx *Expression_listContext) {}

// ExitExpression_list is called when production expression_list is exited.
func (s *BaseModelicaListener) ExitExpression_list(ctx *Expression_listContext) {}

// EnterArray_subscripts is called when production array_subscripts is entered.
func (s *BaseModelicaListener) EnterArray_subscripts(ctx *Array_subscriptsContext) {}

// ExitArray_subscripts is called when production array_subscripts is exited.
func (s *BaseModelicaListener) ExitArray_subscripts(ctx *Array_subscriptsContext) {}

// EnterSubscript is called when production subscript is entered.
func (s *BaseModelicaListener) EnterSubscript(ctx *SubscriptContext) {}

// ExitSubscript is called when production subscript is exited.
func (s *BaseModelicaListener) ExitSubscript(ctx *SubscriptContext) {}

// EnterComment is called when production comment is entered.
func (s *BaseModelicaListener) EnterComment(ctx *CommentContext) {}

// ExitComment is called when production comment is exited.
func (s *BaseModelicaListener) ExitComment(ctx *CommentContext) {}

// EnterString_comment is called when production string_comment is entered.
func (s *BaseModelicaListener) EnterString_comment(ctx *String_commentContext) {}

// ExitString_comment is called when production string_comment is exited.
func (s *BaseModelicaListener) ExitString_comment(ctx *String_commentContext) {}

// EnterAnnotation is called when production annotation is entered.
func (s *BaseModelicaListener) EnterAnnotation(ctx *AnnotationContext) {}

// ExitAnnotation is called when production annotation is exited.
func (s *BaseModelicaListener) ExitAnnotation(ctx *AnnotationContext) {}
