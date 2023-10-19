# Formulas Guide

As introduced in the [Proxeus Handbook](handbook) ([section 5.3](handbook#5.3.4-formulas)), this page describes the most commonly used formulas that will suffice for many users when creating template files.

For more advanced formulas, please refer to the [advanced section](#advanced), or, for the full reference, visit the [Jtwig Book](https://github.com/jtwig/book).

#### IF-Formula

The IF-Formula contains a “condition” and content. This means, if the user enters a certain “Value” (= “condition”), a certain content will be displayed in the document.  
Example: 

`{%if input.Gender == 'male' %} Dear Mr. {%endif%}`

Use “or” to add another “Value” to the condition, displaying the content in both cases.  
Example: 

`{%if input.Gender == 'male' or input.Gender == 'female' %} Dear Mr./Mrs. {%endif%}`

Use “!=” instead of “==” to display the content if the “condition” does not match the “Value”.  
Example: 

`{%if input.Gender != 'male' %} Dear Mrs. {%endif%}`

The same could also be achieved using the IF-ELSE-Formula.

#### IF-ELSE-Formula

The IF-ELSE-Formula displays different content, depending on the “Value” entered.

Example: 

`{%if input.Gender == 'male' %} Dear Mr. {%else%} Dear Mrs. {%endif%}`

#### IF-ELSEIF-ELSE-Formula

Use the IF-ELSEIF-ELSE-Formula to check two “conditions”, displaying different content if the “condition” is matched and showing the content after ELSE when both “conditions” don’t match.

Example: 

`{%if input.Gender == 'male' %} Dear Mr. {%elseif input.Gender == 'female' %} Dear Mrs. {%else%} Dear Mr./Mrs. {%endif%}`

#### FOR Formula

The “FOR” Formula is used to loop through an array of values (i.e. from the Dynamic Lists or Checkboxes) and create a list output.

To gather the birthdays of our friends we have created a field, where multiple entries can be made.

![](handbook/Proxeus%20-%20The%20Complete%20Handbook_html_8a570ce7f4d583b2.png)

To add all the entries from this “Dynamic List” to the document the following formula is used:

```
{% for item in input.FirstName %}
*   {{input.FirstName\[loop.index0\]}} {{input.LastName\[loop.index0\]}}, {{input.Birthdate\[loop.index0\]}}
{% endfor %}
```

This will create a bullet point list (it works with other formatting too):  


![](handbook/Proxeus%20-%20The%20Complete%20Handbook_html_d241ff7211167f3a.png)

For checkboxes, the following formula can be used:

```
{%for item in input.Content %}
• {{item}}
{% endfor %}
```

#### Mathematical Operations

Calculations can be made using the following operators: +, -, \*, /  
Example: Adding one value to another  

`{{input.Value1 + input.Value2}}`

#### Variables

Use the “Set Formula” to define “Values” that are not entered in the form but are needed in your document.

`Use {% set MyValue = 'Content' %} to set a variable.`

Example 1: `{% set MyValue1 = 'This is my set variable one.' %}`

Example 2: `{% set MyValue2 = 'This is my set variable two.' %}`

To set a number for the variable, leave out the ''.

Example 1: `{% set MyNumber1 = 3 %}`

Example 2: `{% set MyNumber2 = 5 %}`

There’s a possibility to set a text list.

Example: `{% set MyList = \['My Value 1', 'My Value 2', 'My Value 3'\] %}`

There are different options to include the variables in the document:

**Input (ODT template file)**

**Output (Document)**

`{{MyValue1}}`

This is my set variable one.

`{{MyValue1 ~ MyValue2}}`

This is my set variable one. This is my set variable two.

`{{MyValue1 ~ ' and ' ~ MyValue2}}`

This is my set variable one. and This is my set variable two.

`{{MyNumber1+MyNumber2}}`

8

`{{MyValue1.contains('is')}}`

true

`{{MyValue1.substring(MyValue1.lastIndexOf('set'))}}`

set variable one.

`{{MyValue1.startsWith('This')}}`

true

`{{MyValue1.endsWith('!')}}`

false

#### Tips

*   IF-Formulas are helpful, if you want to display (or hide) certain text elements in the document, depending on “Values” entered by the user

*   IF-Formulas don’t work if...
    *   they contain umlaute (ä, ö, ü, ...) or special characters (/, &, %, …)
    *   the checkbox or dynamic list component is used (because they allow multiple values)

*   There doesn’t have to be content after `{% else %}`, in this case nothing will be displayed
*   Use `{{myNumber|number\_format(2, '.', ',')}}` to format numbers. The first parameter defines the number of decimal places, the second parameter defines the decimal point and the third parameter defines the thousands separator. The example formula turns 1234 into 1,234.00
*   To check if an optional field has been completed or not (= is empty), use `{% if MyValue1 == '' %} TEXT\_IF\_EMPTY {% else %} TEXT\_IF\_NOT\_EMPTY {% endif %}.` You can also use `{% if MyValue1 != '' %} TEXT\_IF\_NOT\_EMPTY {% endif %}`.

---

<a name="advanced"></a>

# Advanced Formulas

This section describes some useful advanced formulas that might be of interest when creating template files. For the full reference, visit [Jtwig](https://github.com/jtwig/book).

#### Simple Input Output

**Input (ODT template file)**

**Output (Document)**

Current time

`{{'now'|date('dd.MM.yyyy HH:mm:ss.SSSZ')}}`

_See other [](https://docs.oracle.com/javase/7/docs/api/java/text/SimpleDateFormat.html)_ [_date patterns_](https://docs.oracle.com/javase/7/docs/api/java/text/SimpleDateFormat.html)

= `23.08.2019 12:08:33.018+0200.`

Handling null as math value

`{{(null)|number\_format(2, '.', ',')}}`

= 0.00

Calc with text

`{{('2')|number\_format \* ('4')|number\_format}}`

= 8

Custom style for output

```
{% set MyValue3 = 'Value 3' %}

{{MyValue3}}
```

= **Text Value**

Inline Condition

`{{(MyValue3 == 'Value 3')?'MyValue3 equals Value 3':'MyValue3 does not equal Value 3'}}`

= MyValue3 equals Value 3

Replace Text

`{{'I like this and that.'|replace({'this': 'foo', 'that': 'bar'})}}`

= I like foo and bar.

Checks

```
{{defined(\[1,2\]\[5\])}}
= False

{{empty(\[1,2\])}}
= False

{{even(2)}}
= True

{{odd(2)}}
= False

{{iterable(2)}}
= True

{{iterable(\[2,3\])}}
= True

{{first(\[1,2\])}}
= 1

{{last(\[1,2\])}}
= 2
```

Join / Merge

```
{{join(\[1,null,2\],', ')}}
= 1, 2
```

```
{{merge(\[1, 2\], 3)}}
=\[1, 2, 3\]
```

Length

```
{{length(\[1,2\])}}
= 2

{{length(null)}}
= 0

{{length(9)}}
= 1
```

Round

```
{{round(1.33,'CEIL')}}
2

{{round(1.33,'FLOOR')}}
1
```

Capitalize

```
{{capitalize('hello world')}}
Hello world

{{title('hello world')}}
Hello World

{{upper('jtwig')}}
JTWIG

{{lower('jtwig')}}
jtwig
```
