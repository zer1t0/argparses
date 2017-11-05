
// Author: Zer1t0

package argparse

/*
Things to add:
- Basic types
- Maps
- Arrays
*/

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"bytes"
	"path/filepath"
)


const pARAMPREFIX = '-'
const NOSHORTCUT rune = 0

// Types of actions
const (
	actionbegin int = iota
	ActionStoreValue 
	ActionStoreTrue
	ActionStoreFalse
	ActionHelp
	ActionStoreConst
	ActionIncrement
	actionend
)

func isValidAction(action int) bool {
	return action > actionbegin && action < actionend
}


// Arguments categories
const (
	NameCategory int = 1 + iota // Ex: --zeta
	ShortcutCategory // Ex: -z
	ShortcutGroupCategory // Ex: -xyz
	// ParserCategory
	NameValueCategory // Ex: --zeta=1337
	ShortcutValueCategory // Ex: -z=1337
	ShortcutGroupValueCategory // Ex: -xyz=1337
	NameEqCategory // Ex: --zeta=
	ShortcutEqCategory // Ex: -z=
	ShortcutGroupEqCategory // Ex: -xyz=
	ValueCategory // any other value
)



// Util functions
func isLetter(char rune) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

func isValidGroupName(name string) bool {

	if len(name) == 0 {
		return false
	}

	if !isLetter(rune(name[0])){
		return false
	}

	for _, char := range name {
		if !isLetter(char) && char != '-'{
			return false
		}
	}

	return true
}

func isValidName(name string) bool {
	
	var hasLetters = false

	//empty name is correct
	if len(name) == 0 {
		return true
	}

	if len(name) < 3 {
		if !isLetter(rune(name[0])){
			return false
		}
	} else {
		posBegin := name[0] == '-' && name[1] == '-' && isLetter(rune(name[2]))

		opBegin := isLetter(rune(name[0]))

		if !posBegin && !opBegin{
			return false
		} 

	}

	for _, char := range name {
		if isLetter(char){
			hasLetters = true
		}else if char != '-'{
			return false
		}
	}

	return hasLetters
}

func isValidShortcut(char rune) bool {
	return char == 0 || isLetter(char)
}

func isPositionalName(name string) bool {
	
	if len(name) == 0 {
		return false
	}

	if len(name) < 3 {
		return true
	}

	if name[0] == '-' && name[1] == '-' {
		return false
	}

	return true
}




// ==INTERFACE value==
type value interface{
    get() string
    set(string, string, rune) error
    setDefault()
    setTrue() error
    setFalse() error
    setConstant() error
    increment() error
}

// ==CLASS intValue BEGIN==
type IntArgCheckFunc func(int) bool

type intValue struct {
    constValue   int
    defaultValue int
    value        *int
    checkValue	IntArgCheckFunc
}

func newIntValue(defaultValue int, constValue int, checkValue IntArgCheckFunc) intValue{
    val := intValue{}
    val.defaultValue = defaultValue
    val.constValue = constValue
    val.value = new(int)
    val.checkValue = checkValue

    val.setDefault()

    return val
}

func (val intValue) get() string {
	return fmt.Sprint(*(val.value))
}

func (val intValue) set(value string, name string, shortcut rune) error {
    vali, err := strconv.ParseInt(value, 0, 0)
    if err != nil {
        return fmt.Errorf("Invalid value \"%s\" for argument %s[-%c], must be an integer", value, name, shortcut)
    }
    valInt := int(vali)

    if val.checkValue != nil{
    	if !val.checkValue(valInt){
    		return fmt.Errorf("Invalid value \"%s\" for argument %s[-%c], must meet custom restriction", value, name, shortcut)
    	}
    }

    *(val.value) = valInt
    return nil
}

func (val intValue) setDefault() {
	*(val.value) = val.defaultValue
}

func (val intValue) setTrue() error {
    return fmt.Errorf("Invalid action for int: Store True")
}

func (val intValue) setFalse() error {
    return fmt.Errorf("Invalid action for int: Store False")
}

func (val intValue) setConstant() error {
    *(val.value) = val.constValue
    return nil
}

func (val intValue) increment() error{
	*(val.value) += 1
    return nil
}
// ==CLASS boolValue END==


// ==CLASS stringValue BEGIN==
type StringArgCheckFunc func(string) bool

type stringValue struct{
    constValue   string
    defaultValue string
    value        *string
    checkValue StringArgCheckFunc
}

func newStringValue(defaultValue string, constValue string, checkValue StringArgCheckFunc) stringValue{
    val := stringValue{}
    val.defaultValue = defaultValue
    val.constValue = constValue
    val.value = new(string)
    val.checkValue = checkValue
    
    val.setDefault()

    return val
}

func (arg stringValue) get() string {
	return *(arg.value)
}

func (val stringValue) set(value string, name string, shortcut rune) error {
    
    if val.checkValue != nil{
    	if !val.checkValue(value){
    		return fmt.Errorf("Invalid value \"%s\" for argument %s[-%c], must meet custom restriction", value, name, shortcut)
    	}
    }

    *(val.value) = value
    return nil
}

func (val stringValue) setDefault() {
	*(val.value) = val.defaultValue
}

func (val stringValue) setTrue() error {
    return fmt.Errorf("Invalid action for string: Store True")
}

func (val stringValue) setFalse() error {
    return fmt.Errorf("Invalid action for int: Store False")
}


func (val stringValue) setConstant() error {
    *(val.value) = val.constValue
    return nil
}

func (val stringValue) increment() error{
    return fmt.Errorf("Invalid action for string: Increment")
}
// ==CLASS boolValue END==


// ==CLASS boolValue BEGIN==
type BoolArgCheckFunc func(bool) bool

type boolValue struct {
	constValue   bool
	defaultValue bool
	value        *bool
	checkValue	BoolArgCheckFunc
}

func newBoolValue(defaultValue bool, constValue bool, checkValue BoolArgCheckFunc) boolValue{
    val := boolValue{}
    val.constValue = constValue
    val.defaultValue = defaultValue
    val.value = new(bool)
    val.checkValue = checkValue

    val.setDefault()

    return val
}

func (val boolValue) get() string { 
	return fmt.Sprint(*(val.value))
}

func (val boolValue) set(value string, name string, shortcut rune) error {
	valb, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("Invalid value \"%s\" for argument %s[-%c], must be a boolean", value, name, shortcut)
	}

	valBool := bool(valb)
	if val.checkValue != nil{
    	if !val.checkValue(valBool){
    		return fmt.Errorf("Invalid value \"%s\" for argument %s[-%c], must meet custom restriction", value, name, shortcut)
    	}
    }

	*(val.value) = valBool
	return nil
}

func (val boolValue) setDefault() {
	*(val.value) = val.defaultValue
}

func (val boolValue) setTrue() error {
    *(val.value) = true
    return nil
}

func (val boolValue) setFalse() error {
    *(val.value) = false
    return nil
}


func (val boolValue) setConstant() error {
    *(val.value) = val.constValue
    return nil
}

func (val boolValue) increment() error{
    return fmt.Errorf("Invalid action for bool: Increment")
}
// ==CLASS boolValue END==


// ==CLASS argument BEGIN==
type argument struct{
    name         string
    shortcut     rune
    description  string
    action       int
    mandatory    bool
    positional   bool
    val          value
}

// CONSTRUCTORS of argument
func  newArgument(name string, shortcut rune, description string, mandatory bool) *argument {
	
	var arg = new(argument)

	arg.name = strings.ToLower(name)
	arg.shortcut = shortcut
	arg.description = description
	arg.action = ActionStoreValue
	arg.mandatory = mandatory
	arg.positional = isPositionalName(name)

	return arg
}


// PUBLIC METHODS OF arguments
func (arg argument) usage() string {

	if arg.positional {
		return fmt.Sprintf("%s", arg.name)
	}

	if arg.name != ""{
		if arg.shortcut != NOSHORTCUT {
			if arg.action == ActionStoreValue {
				return fmt.Sprintf("%s/-%c %s", arg.name, arg.shortcut, strings.ToUpper(arg.name[2:]))
			} else {
				return fmt.Sprintf("%s/-%c", arg.name, arg.shortcut)
			}
		} else {
			if arg.action == ActionStoreValue {
				return fmt.Sprintf("%s %s", arg.name, strings.ToUpper(arg.name[2:]))
			} else {
				return fmt.Sprintf("%s", arg.name)
			}
		}
	} else {
		if arg.action == ActionStoreValue {
			return fmt.Sprintf("-%c %s", arg.shortcut, strings.ToUpper(string(arg.shortcut)))
		} else {
			return fmt.Sprintf("-%c", arg.shortcut)
		}
	}
}

func (arg argument) help() string{
	return fmt.Sprintf("%s\t%s", arg.usage(), arg.description)
}

func (arg argument) set(value string) error {
    return arg.val.set(value, arg.name, arg.shortcut)
}

func (arg argument) setDefault() {
	arg.val.setDefault()
}

func (arg argument) get() string {
    return arg.val.get()
}

func (arg argument) setTrue() error {
    return arg.val.setTrue()
}

func (arg argument) setFalse() error {
    return arg.val.setFalse()
}

func (arg argument) setConstant() error {
    return arg.val.setConstant()
}

func (arg argument) increment() error{
    return arg.val.increment()
}

// ==CLASS argument END==




// ==CLASS argumentsGroup==
type argumentsGroup struct {
	name      string
	description string
	parser    *argParser
	arguments []*argument
	required  bool
	exclusive bool
}


// CONSTRUCTORS of argumentsGroup
func newArgumentGroup(name string, description string, required bool, exclusive bool, parser *argParser) *argumentsGroup{
	group := new(argumentsGroup)

	group.name = strings.ToLower(name)
	group.description = description
	group.parser = parser
	group.arguments = make([]*argument, 0, 8)
	group.required = required
	group.exclusive = exclusive

	return group
}


// PUBLIC METHODS of argumentsGroup
func (argGroup *argumentsGroup) AddInt(name string, shortcut rune, description string, mandatory bool, action int, defaultValue int, constValue int, checkValue IntArgCheckFunc) (*int, error){
	return argGroup.parser.AddInt(name, shortcut, description, mandatory, action, defaultValue, constValue, checkValue, argGroup.name)
}

func (argGroup *argumentsGroup) AddString(name string, shortcut rune, description string, mandatory bool, action int, defaultValue string, constValue string, checkValue StringArgCheckFunc) (*string, error){
	return argGroup.parser.AddString(name, shortcut, description, mandatory, action, defaultValue, constValue, checkValue, argGroup.name)
}

func (argGroup *argumentsGroup) AddBool(name string, shortcut rune, description string, mandatory bool, action int, defaultValue bool, constValue bool, checkValue BoolArgCheckFunc) (*bool, error){
	return argGroup.parser.AddBool(name, shortcut, description, mandatory, action, defaultValue, constValue, checkValue, argGroup.name)
}

// PRIVATE METHODS of argumentsGroup
func (argGroup *argumentsGroup) addArgument(arg *argument) {
	argGroup.arguments = append(argGroup.arguments, arg)
}

func (argGroup *argumentsGroup) getArgument(name string) (*argument, bool){
	if name == "" {
		return nil, false
	}

	lowerName := strings.ToLower(name)

	for _, arg := range argGroup.arguments {
		if lowerName == arg.name {
			return arg, true
		}
	}
	return nil, false
}

func (argGroup *argumentsGroup) existsArgument(name string) bool {
	_, exists := argGroup.getArgument(name)
	return exists
}


func (argGroup *argumentsGroup) usage() string{

	usages := make([]string, 0, 16)


	if argGroup.exclusive {
		for _, arg := range argGroup.arguments {
				usages = append(usages, arg.usage())
		}

	
		return fmt.Sprintf("[%s]", strings.Join(usages , " | "))
	}

	for _, arg := range argGroup.arguments {
				usages = append(usages, fmt.Sprintf("[%s]", arg.usage()))
	}

	return strings.Join(usages , " ")
}

// ==CLASS argumentsGroup END==





// ==CLASS argParser BEGIN==
type argParser struct {
	name       	string
	description string
	prefix     	rune
	helpArgument *string
	arguments  	[]*argument
	posArguments []*argument
	groups     map[string]*argumentsGroup
	subparsers map[string]*argParser
	subparserRequired bool
	selectedSubparser *string
}


// CONSTRUCTORS of argParser
func NewArgParser(name string, description string, includeHelp bool) (*argParser, error) {
	parser := new(argParser)

	if name == "" {
		parser.name = filepath.Base(os.Args[0])
	} else {
		if !isValidGroupName(name) {
			return nil, fmt.Errorf("Invalid name %s for parser, it must begin for letter and only contains letters and -", name)
		}
		parser.name = strings.ToLower(name)
	}

	parser.description =  description
	parser.arguments = []*argument{}
	parser.posArguments = []*argument{}
	parser.groups = map[string]*argumentsGroup{}
	parser.subparsers = map[string]*argParser{}
	parser.subparserRequired = false
	parser.prefix = pARAMPREFIX
	parser.selectedSubparser = new(string)

	// adds help param
	
	if includeHelp {
		aux, err := parser.AddString("--help", 'h', "Print this message", false, ActionHelp, "", "", nil, "")
		
		if err != nil {
			fmt.Println(err)
		} else {
			parser.helpArgument = aux
		}
	}

	return parser, nil
}


// PRIVATE METHODS of argParser
func (parser *argParser) existsGroup(name string) bool {
	_, exists := parser.groups[strings.ToLower(name)]
	return exists
}

func (parser *argParser) existsSubparser(name string) bool {
	_, exists := parser.subparsers[strings.ToLower(name)]
	return exists
}

func (parser *argParser) existsArgument(name string) bool {
	_, exists := parser.getArgument(name)
	return exists
}

func (parser *argParser) existsArgumentInGroups(name string) bool{
	
	if name == "" {
		return false
	}

	lowerName := strings.ToLower(name)

	for _, grp := range parser.groups {
		if grp.existsArgument(lowerName) {
			return true
		}
	}

	return false
}

func (parser *argParser) existsArgumentFromShortcut(shortcut rune) bool {
	_, exists := parser.getArgumentFromShortcut(shortcut)
	return exists
}

func (parser *argParser) getArgumentCategory(argument string) int{

	var argLen = len(argument)
	var category = ValueCategory
	//var char rune = 0

	if argLen == 0 {
		return ValueCategory
	}


	if argument[0] == '-'{
		if argLen == 1 {
			return ValueCategory
		}

		if argument[1] == '-'{
			if argLen == 2 {
				return ValueCategory
			} else {
				//could be an argument name
				groups := strings.Split(argument, "=")

				name := groups[0]

				_, exists := parser.getArgument(name)

				if !exists{
					return ValueCategory
				}

				if len(groups) == 1 {
					return NameCategory
				}

				if len(groups) > 2 {
					return NameValueCategory
				}

				if len(groups[1]) == 0 {
					return NameEqCategory
				} else {
					return NameValueCategory
				}
				
			}
		} else {
			// could be a Shortcut
			_, exists := parser.getArgumentFromShortcut(rune(argument[1]))
			
			if !exists {
				return ValueCategory
			} else {
				category = ShortcutCategory
			}

			for i := 2; i < len(argument); i++{
				_, exists = parser.getArgumentFromShortcut(rune(argument[i]))

				if !exists {
					if argument[i] == '='{
						if category == ShortcutCategory {
							if i == len(argument) - 1{
								return ShortcutEqCategory
							} else {
								return ShortcutValueCategory
							}
						} else if category == ShortcutGroupCategory {
							if i == len(argument) - 1{
								return ShortcutGroupEqCategory
							} else {
								return ShortcutGroupValueCategory
							}
						} else {
							fmt.Printf("Oh oh, you shouldn't see this (Group category fail)\n")
						}
					} else {
						return ValueCategory
					}
				}

				category = ShortcutGroupCategory	
			}

			return category
		}

	}

	return category
}

func (parser *argParser) getArgumentFromShortcut(shortcut rune) (*argument, bool) {

	if shortcut == NOSHORTCUT {
		return nil, false
	}

	for _, arg := range parser.arguments {
		if shortcut == arg.shortcut {
			return arg, true
		}
	}
	return nil, false
}

func (parser *argParser) getArgument(name string) (*argument, bool){
	if name == "" {
		return nil, false
	}

	lowerName := strings.ToLower(name)

	for _, arg := range parser.arguments {
		if lowerName == arg.name {
			return arg, true
		}
	}

	//arg, ok := parser.arguments[strings.ToLower(name)]
	return nil, false
}

func (parser *argParser) getGroup(name string) (*argumentsGroup, bool){
	group, ok := parser.groups[strings.ToLower(name)]
	return group, ok
}

func (parser *argParser) getSubparser(name string) (*argParser, bool){
	subparser, ok := parser.subparsers[strings.ToLower(name)]
	return subparser, ok
}

func (parser *argParser) setDefaultValues() {

	for _, arg := range parser.arguments {
		arg.setDefault()
	}

}

func (parser *argParser) isValidFlagName(name string) bool {

	//empty name is correct
	if len(name) == 0 {
		return true
	}

	if len(name) < 3 {
		//name must start with letter
		if !isLetter(rune(name[0])){
			return false
		}
	} else {
		//name must start with letter or -- 
		posBegin := rune(name[0]) == parser.prefix && rune(name[1]) == parser.prefix && isLetter(rune(name[2]))

		opBegin := isLetter(rune(name[0]))

		if !posBegin && !opBegin{
			return false
		} 

	}

	//only letters and - are available for flags names
	for _, char := range name {
		if !isLetter(char) && char != '-'{
			return false
		}
	}

	return true
}

func (parser *argParser) createArg(name string, shortcut rune, description string, mandatory bool) (*argument, error) {

	var arg = newArgument(name, shortcut, description, mandatory)

	if name == "" && shortcut == NOSHORTCUT {
		return nil, fmt.Errorf("Argument must have a name and/or shortcut")
	}

	if !parser.isValidFlagName(arg.name) {
		return nil, fmt.Errorf("Argument name only can be letters and -")
	}

	if parser.existsArgument(arg.name) {
		return nil, fmt.Errorf("Argument %s is already defined in parser %s", name, parser.name)
	}


	if arg.positional {
		arg.mandatory = true
		arg.shortcut = NOSHORTCUT

	} else {
		// check shortcut
		if shortcut == NOSHORTCUT {
			arg.shortcut = NOSHORTCUT
		} else {
			if !isValidShortcut(shortcut) {
				return nil, fmt.Errorf("Invalid shorcut, must be a letter")
			}

			argEx, exists := parser.getArgumentFromShortcut(shortcut)
			if exists {
				return nil, fmt.Errorf("Shortcut %c (arg %s) is already used in argument %s (parser %s)", shortcut, name, argEx.name, parser.name)
			}
		}

	}

	return arg, nil
}

func (parser *argParser) addArg(arg *argument, group string) error {
	if group != "" {
		if !parser.existsGroup(group) {
			return fmt.Errorf("Group %s is not defined", group)
		}
		if arg.positional {
			return fmt.Errorf("Positional argument %s cannot belong to a group", arg.name)
		}
		parser.groups[group].addArgument(arg)
	}

	parser.arguments = append(parser.arguments, arg)
	//parser.arguments[arg.name] = arg

	if arg.positional {
		parser.posArguments = append(parser.posArguments, arg)
	}

	return nil
}

// PUBLIC METHODS OF argParser

func (parser *argParser) SetFlagPrefix(prefix rune) {
	parser.prefix = prefix
}

func (parser *argParser) SetSubparserRequired(required bool){
	parser.subparserRequired = required
}

func (parser *argParser) GetSelectedSubparser() string {
	return *(parser.selectedSubparser)
}

func (parser *argParser) GetHelpArgument() *string {
	return parser.helpArgument
}

func (parser *argParser) AddInt(name string, shortcut rune, description string, mandatory bool, action int, defaultValue int, constValue int, checkValue IntArgCheckFunc, group string) (*int, error) {

	arg, err := parser.createArg(name, shortcut, description, mandatory)
	
	if err != nil {
		return nil, err
	}

	val := newIntValue(defaultValue, constValue, checkValue)

	if arg.mandatory {
		arg.action = ActionStoreValue
	} else {

		if isValidAction(action) {
			arg.action = action
		} else {
			return nil, fmt.Errorf("Invalid action")
		}

		switch action {
			case ActionStoreTrue:
				return nil, fmt.Errorf("Store True is only available in booleans")

			case ActionStoreFalse:
				return nil, fmt.Errorf("Store False is only available in booleans")

			case ActionHelp:
				return nil, fmt.Errorf("Help is only available in strings")
		}
	}

	arg.val = val

	err = parser.addArg(arg, group)

	if err != nil {
		return nil, err
	}

	return val.value, nil
}

func (parser *argParser) AddString(name string, shortcut rune, description string, mandatory bool, action int, defaultValue string, constValue string, checkValue StringArgCheckFunc, group string) (*string, error){
	
	arg, err := parser.createArg(name, shortcut, description, mandatory)

	if err != nil {
		return nil, err
	}

	val := newStringValue(defaultValue, constValue, checkValue)

	if arg.mandatory {
		arg.action = ActionStoreValue
	} else {

		if isValidAction(action) {
			arg.action = action
		} else {
			return nil, fmt.Errorf("Invalid action")
		}

		switch action {
			case ActionStoreTrue:
				return nil, fmt.Errorf("Store True is only available in booleans")

			case ActionStoreFalse:
				return nil, fmt.Errorf("Store False is only available in booleans")

			case ActionIncrement:
				return nil, fmt.Errorf("Increment is only available in integers")
		}
	}

	arg.val = val

	err = parser.addArg(arg, group)

	if err != nil {
		return nil, err
	}

	return val.value, nil
}

func (parser *argParser) AddBool(name string, shortcut rune, description string, mandatory bool, action int, defaultValue bool, constValue bool, checkValue BoolArgCheckFunc, group string) (*bool, error) {

	arg, err := parser.createArg(name, shortcut, description, mandatory)

	if err != nil {
		return nil, err
	}

	val := newBoolValue(defaultValue, constValue, checkValue)

	if arg.mandatory {
		arg.action = ActionStoreValue
	} else {

		if isValidAction(action) {
			arg.action = action
		} else {
			return nil, fmt.Errorf("Invalid action")
		}

		switch action {
			case ActionStoreTrue:
				val.defaultValue = false

			case ActionStoreFalse:
				val.defaultValue = true

			case ActionIncrement:
				return nil, fmt.Errorf("Increment is only available in integers")

			case ActionHelp:
				return nil, fmt.Errorf("Help is only available in strings")
		}

	}

	arg.val = val

	err = parser.addArg(arg, group)

	if err != nil {
		return nil, err
	}

	return val.value, nil
}

func (parser *argParser) AddArgumentsGroup(name string, description string, required bool, exclusive bool) (*argumentsGroup, error) {
	
	var argsGroup = newArgumentGroup(name, description, required, exclusive, parser)

	if !isValidGroupName(argsGroup.name) {
		return nil, fmt.Errorf("Invalid name %s for group, it must begin for letter and only contains letters and -", argsGroup.name)
	}

	if parser.existsGroup(argsGroup.name) {
		return nil, fmt.Errorf("Group %s is already defined", argsGroup.name)
	}

	parser.groups[argsGroup.name] = argsGroup

	return argsGroup, nil
}

func (parser *argParser) AddSubparser(name string, description string, includeHelp bool) (*argParser, error) {

	if name == "" {
		return nil, fmt.Errorf("Name for subparser cannot be empty")
	}

	if parser.existsSubparser(name) {
		return nil, fmt.Errorf("Subparser %s of %s is already defined", name, parser.name)
	}

	subparser, err := NewArgParser(name, description, includeHelp)

	if err != nil {
		return nil, err
	}

	parser.subparsers[strings.ToLower(name)] = subparser

	return subparser, nil
}



func (parser *argParser) Usage() string {

	var usage bytes.Buffer

	usage.WriteString(fmt.Sprintf("Usage: %s ", parser.name))

	
	for _, grp := range parser.groups {
		usage.WriteString(fmt.Sprintf("%s ",grp.usage()))
	}


	for _, arg := range parser.arguments {
		if !arg.positional && !parser.existsArgumentInGroups(arg.name){
			usage.WriteString(fmt.Sprintf("[%s] ", arg.usage()))
		}
	}

	for _, arg := range parser.posArguments {
		usage.WriteString(fmt.Sprintf("%s ", arg.usage()))
	}

	if len(parser.subparsers) > 0 {
		parserNames := []string{}

		for _, subparser := range parser.subparsers{
			parserNames = append(parserNames, subparser.name)
		}

		usage.WriteString(fmt.Sprintf("{%s} ...", strings.Join(parserNames, ",")))
	}

	return usage.String()
}

func (parser *argParser) Help() string {

	var help bytes.Buffer

	help.WriteString(parser.Usage()+"\n")


	if parser.description != "" {
		help.WriteString("\n" + parser.description + "\n")
	}

	if len(parser.posArguments) > 0 || len(parser.subparsers) > 0 {
		help.WriteString("\nPositional arguments:\n")
	}

	for _, arg := range parser.posArguments {
		help.WriteString(arg.help() + "\n")
	}

	if len(parser.subparsers) > 0 {
		parserNames := []string{}

		for _, subparser := range parser.subparsers{
			parserNames = append(parserNames, subparser.name)
		}

		help.WriteString(fmt.Sprintf("{%s}\tsubcommands\n", strings.Join(parserNames, ",")))
	}

	if len(parser.arguments) > 0 {
		help.WriteString("\nOptional arguments:\n")
	}

	for _, arg := range parser.arguments {
		if !arg.positional{
			help.WriteString(fmt.Sprintf("%s\n", arg.help()))
		}
	}

	return help.String()
}

// Parse parse arguments and return the parameters
func (parser *argParser) Parse(arguments []string) (error) {
	
	argStr := ""
	var errMessage bytes.Buffer
	currentArgs := arguments
	index := 0
	var currentArg *argument = nil
	argsSet := map[string]*argument{}
	positionalIndex := 0
	numNonProcessesPositionals := 0

	*(parser.selectedSubparser) = ""

	parser.setDefaultValues()

	if arguments == nil {
		currentArgs = os.Args
	}

	if len(currentArgs) < 1 {
		return fmt.Errorf("No command was provided")
	}

	currentArgs = currentArgs[1:]
	numNonProcessesPositionals = len(parser.posArguments)

	ParseLoop: for index = 0; index < len(currentArgs) - numNonProcessesPositionals; index++{
		argStr = currentArgs[index]
		category := parser.getArgumentCategory(argStr)

		switch category {
			case ValueCategory:
				
				if positionalIndex < len(parser.posArguments){
					err := parser.posArguments[positionalIndex].set(argStr)
					if err != nil{
						return err
					}
					positionalIndex++
					numNonProcessesPositionals--
				} else {
					break ParseLoop
				}
				continue

			case NameCategory:
				currentArg, _ = parser.getArgument(argStr)

				switch currentArg.action{
					case ActionStoreValue:
						if index < len(currentArgs) - 1 {
							index++
							err := currentArg.set(currentArgs[index])
							if err != nil {
								return err
							}
						} else {
							return fmt.Errorf("No value for argument %s", currentArg.name)
						}
					case ActionStoreTrue:
						currentArg.setTrue()

					case ActionStoreFalse:
						currentArg.setFalse()
					
					case ActionIncrement:
						currentArg.increment()

					case ActionStoreConst:
						currentArg.setConstant()

					case ActionHelp:
						currentArg.set(parser.Help())
						return fmt.Errorf("Help")

					default:
						fmt.Printf("This is a top secret message, or maybe a bug\n")

				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue
				
			case ShortcutCategory:
				currentArg, _ = parser.getArgumentFromShortcut(rune(argStr[1]))

				switch currentArg.action{
					case ActionStoreValue:
						if index < len(currentArgs) - 1 {
							index++
							err := currentArg.set(currentArgs[index])
							if err != nil {
								return err
							}
						} else {
							return fmt.Errorf("No value for argument %s", currentArg.name)
						}
					case ActionStoreTrue:
						currentArg.setTrue()

					case ActionStoreFalse:
						currentArg.setFalse()
					
					case ActionIncrement:
						currentArg.increment()

					case ActionStoreConst:
						currentArg.setConstant()

					case ActionHelp:
						currentArg.set(parser.Help())
						return fmt.Errorf("Help")

					default:
						fmt.Printf("This is a top secret message, or maybe a bug\n")

				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue

			case ShortcutGroupCategory:

				for i := 1; i < len(argStr); i++ {
					currentArg, _ = parser.getArgumentFromShortcut(rune(argStr[i]))

					switch currentArg.action{
						case ActionStoreValue:
							if i < len(argStr) - 1 {
								err := currentArg.set(argStr[(i+1):])
								if err != nil {
									return err
								}
							} else if index < len(currentArgs) - 1 {
								index++
								err := currentArg.set(currentArgs[index])
								if err != nil {
									return err
								}

							} else {
								return fmt.Errorf("No value for argument %s", currentArg.name)
							}
							break
						case ActionStoreTrue:
							currentArg.setTrue()

						case ActionStoreFalse:
							currentArg.setFalse()
						
						case ActionIncrement:
							currentArg.increment()

						case ActionStoreConst:
							currentArg.setConstant()

						case ActionHelp:
							currentArg.set(parser.Help())
							return fmt.Errorf("Help")

						default:
							fmt.Printf("This is a top secret message, or maybe a bug\n")
					}
				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue

			case NameValueCategory:
				groups := strings.Split(argStr, "=")

				argName := groups[0]
				argValue := strings.Join(groups[1:], "=")

				currentArg, _ = parser.getArgument(argName)

				switch currentArg.action{
					case ActionStoreValue:
						err := currentArg.set(argValue)
						if err != nil {
							return err
						}
					case ActionStoreTrue:
						currentArg.setTrue()

					case ActionStoreFalse:
						currentArg.setFalse()
					
					case ActionIncrement:
						currentArg.increment()

					case ActionStoreConst:
						currentArg.setConstant()

					case ActionHelp:
						currentArg.set(parser.Help())
						return fmt.Errorf("Help")

					default:
						fmt.Printf("This is a top secret message, or maybe a bug\n")

				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue

			case ShortcutValueCategory:
				groups := strings.Split(argStr, "=")

				argShortcut := groups[0][1]
				argValue := strings.Join(groups[1:], "=")

				currentArg, _ = parser.getArgumentFromShortcut(rune(argShortcut))

				switch currentArg.action {
					case ActionStoreValue:
						err := currentArg.set(argValue)
						if err != nil {
							return err
						}
					case ActionStoreTrue:
						currentArg.setTrue()

					case ActionStoreFalse:
						currentArg.setFalse()
					
					case ActionIncrement:
						currentArg.increment()

					case ActionStoreConst:
						currentArg.setConstant()

					case ActionHelp:
						currentArg.set(parser.Help())
						return fmt.Errorf("Help")

					default:
						fmt.Printf("This is a top secret message, or maybe a bug\n")

				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue

			case ShortcutGroupValueCategory:

				groups := strings.Split(argStr, "=")
				groupShortcuts := groups[0]
				argValue := strings.Join(groups[1:], "=")

				for i := 1; rune(argStr[i]) != '='; i++ {
					currentArg, _ = parser.getArgumentFromShortcut(rune(argStr[i]))

					switch currentArg.action{
						case ActionStoreValue:
							if i < len(groupShortcuts) - 1 {
								err := currentArg.set(argStr[(i+1):])
								if err != nil {
									return err
								}
							} else {
								err := currentArg.set(argValue)
								if err != nil {
									return err
								}
							}
							break
						case ActionStoreTrue:
							currentArg.setTrue()

						case ActionStoreFalse:
							currentArg.setFalse()
						
						case ActionIncrement:
							currentArg.increment()

						case ActionStoreConst:
							currentArg.setConstant()

						case ActionHelp:
							currentArg.set(parser.Help())
							return fmt.Errorf("Help")

						default:
							fmt.Printf("This is a top secret message, or maybe a bug\n")
					}
					argsSet[currentArg.name] = currentArg
				}
				currentArg = nil
				continue

			case NameEqCategory:
				groups := strings.Split(argStr, "=")
				argName := groups[0]
				currentArg, _ = parser.getArgument(argName)

				switch currentArg.action{
					case ActionStoreValue:
						if index < len(currentArgs) - 1 {
							index++
							err := currentArg.set(currentArgs[index])
							if err != nil {
								return err
							}
						} else {
							return fmt.Errorf("No value for argument %s", currentArg.name)
						}
					case ActionStoreTrue:
						currentArg.setTrue()

					case ActionStoreFalse:
						currentArg.setFalse()
					
					case ActionIncrement:
						currentArg.increment()

					case ActionStoreConst:
						currentArg.setConstant()

					case ActionHelp:
						currentArg.set(parser.Help())
						return fmt.Errorf("Help")

					default:
						fmt.Printf("This is a top secret message, or maybe a bug\n")

				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue

			case ShortcutEqCategory:
				groups := strings.Split(argStr, "=")
				argShortcut := groups[0][1]
				currentArg, _ = parser.getArgumentFromShortcut(rune(argShortcut))

				switch currentArg.action{
					case ActionStoreValue:
						if index < len(currentArgs) - 1 {
							index++
							err := currentArg.set(currentArgs[index])
							if err != nil {
								return err
							}
						} else {
							return fmt.Errorf("No value for argument %s", currentArg.name)
						}
					case ActionStoreTrue:
						currentArg.setTrue()

					case ActionStoreFalse:
						currentArg.setFalse()
					
					case ActionIncrement:
						currentArg.increment()

					case ActionStoreConst:
						currentArg.setConstant()

					case ActionHelp:
						currentArg.set(parser.Help())
						return fmt.Errorf("Help")

					default:
						fmt.Printf("This is a top secret message, or maybe a bug\n")

				}
				argsSet[currentArg.name] = currentArg
				currentArg = nil
				continue
			
			case ShortcutGroupEqCategory:
				groups := strings.Split(argStr, "=")
				groupShortcuts := groups[0]

				for i := 1; rune(argStr[i]) != '='; i++ {
					currentArg, _ = parser.getArgumentFromShortcut(rune(argStr[i]))

					switch currentArg.action{
						case ActionStoreValue:
							if i < len(groupShortcuts) - 1 {
								err := currentArg.set(argStr[(i+1):])
								if err != nil {
									return err
								}
							} else if index < len(currentArgs) - 1 {
								index++
								err := currentArg.set(currentArgs[index])
								if err != nil {
									return err
								}
							} else {
								return fmt.Errorf("No value for argument %s", currentArg.name)
							}
							break
						case ActionStoreTrue:
							currentArg.setTrue()

						case ActionStoreFalse:
							currentArg.setFalse()
						
						case ActionIncrement:
							currentArg.increment()

						case ActionStoreConst:
							currentArg.setConstant()

						case ActionHelp:
							currentArg.set(parser.Help())
							return fmt.Errorf("Help")

						default:
							fmt.Printf("This is a top secret message, or maybe a bug\n")
					}
					argsSet[currentArg.name] = currentArg
				}
				currentArg = nil
				continue
			
			default:
				fmt.Printf("What? Are you seeing me? I'm a fail in code (switch category)\n")
				return fmt.Errorf("Fail in code")
		}
	}

	//if only there are one argument, check if it is a help argument
	if len(currentArgs) == 1 && numNonProcessesPositionals == 1 {
		argStr = currentArgs[0]
		category := parser.getArgumentCategory(argStr)

		if category == NameCategory{
			currentArg, _ = parser.getArgument(argStr)

			if currentArg.action == ActionHelp {
				currentArg.set(parser.Help())
				return fmt.Errorf("Help")
			}
						
		} else if category == ShortcutCategory {

			currentArg, _ = parser.getArgumentFromShortcut(rune(argStr[1]))

			if currentArg.action == ActionHelp {
				currentArg.set(parser.Help())
				return fmt.Errorf("Help")
			}

		}
	}


	// process the rest of positional arguments
	for ; positionalIndex < numNonProcessesPositionals && index < len(currentArgs); positionalIndex++{

		err := parser.posArguments[positionalIndex].set(currentArgs[index])
		if err != nil{
			return err
		}
		index++
	}

	// check if positional arguments were set
	if positionalIndex < len(parser.posArguments) {
		return fmt.Errorf("Too few arguments were provided")
	}

	// check if mandatory arguments were set
	for _, arg := range parser.arguments {
		if !arg.mandatory || arg.positional {continue}

		_, wasSet := argsSet[arg.name]
		if !wasSet{
			errMessage.WriteString(fmt.Sprintf("argument %s has no value\n", arg.name))
		}
	}

	/*if errMessage.Len() > 0 {
		return nil, fmt.Errorf("%s", errMessage.String())
	}*/


	// GROUPS
	// check requirements of the arguments groups
	for _, argGroup := range parser.groups {

		argGroupSet := make([]*argument, 0, 8)

		for _, arg := range argGroup.arguments {
			_, wasSet := argsSet[arg.name]
			if wasSet{
				argGroupSet = append(argGroupSet, arg)
			}

		}

		if len(argGroupSet) == 0 && argGroup.required {
			
			argNames := make([]string,0,8)
			for _, arg := range argGroup.arguments {
				argNames = append(argNames, fmt.Sprintf("%s[-%c]",arg.name, arg.shortcut))
			}
			errMessage.WriteString(fmt.Sprintf("no argument from group \"%s\" (%s) was specified\n", argGroup.name, strings.Join(argNames,",")))
	
		} else if len(argGroupSet) > 1 && argGroup.exclusive {
			argNames := make([]string,0,8)
			for _, arg := range argGroupSet {
				argNames = append(argNames, fmt.Sprintf("%s[-%c]",arg.name, arg.shortcut))
			}

			errMessage.WriteString(fmt.Sprintf("more than one argument of group \"%s\" was specified (%s)\n", argGroup.name, strings.Join(argNames,",")))
		}
	}

	if errMessage.Len() > 0 {
		return fmt.Errorf("%s", errMessage.String())
	}

	// check if there are enough arguments for subparser
	if index < len(currentArgs) {
		subparser, ok := parser.getSubparser(currentArgs[index])

		if ok {
			*(parser.selectedSubparser) = subparser.name
			// parse subparser
			err := subparser.Parse(currentArgs[index:])

			if fmt.Sprintf("%s",err) == "Help"{
				help := subparser.GetHelpArgument()
				*(parser.helpArgument) = *help
				return err
			}

		}

	}

	if parser.subparserRequired {
		names := make([]string,0,8)
		for _, subparser := range parser.subparsers {
			names = append(names, subparser.name)
		}
		return fmt.Errorf("no subcommand was specified (%s)\n", strings.Join(names,","))
	}


	return nil
}

// ==CLASS argParser END==


// function to split a string into different arguments
func StringToArgv(line string) []string {

    args := make([]string, 0, 32)
    finisher := ' '
    escaped := false
    arg := bytes.Buffer{}
    i := 0

    subline := strings.TrimSpace(line)

    fmt.Println(subline)

    for len(subline) > 0 {
        arg.Reset()
        finisher = ' '
        escaped = false

        switch (subline[0]){
            case '"':
                finisher = '"'
            case '\'':
                finisher = '\''
            case '\\':
                escaped = true
            default:
                arg.WriteRune(rune(subline[0]))
        }

        for i = 1; i < len(subline) && !(rune(subline[i]) == finisher && !escaped); i++ {

            if subline[i] == '\\' && !escaped {
                escaped = true
            } else {
                arg.WriteRune(rune(subline[i]))
                escaped = false
            }
        }
        args = append(args, arg.String())
        i++
        if i >= len(subline) {
            subline = ""
        } else {
            subline = strings.TrimSpace(subline[i:])
        }
    }

    return args
}
