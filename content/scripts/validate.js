function validateFood() {
	var result = false;

	var resName = validateRequiredString("name");
	var resCal = validateRequiredNumeric("calories");
	var resPro = validateRequiredNumeric("protein");
	var resCar = validateRequiredNumeric("carbohydrate");
	var resFat = validateRequiredNumeric("fat");

	var result = (resName && resCal && resPro && resCar && resFat);

	return result; 
}

function addSuccessToInput(elementName){
	$("[name='"+ elementName + "']").removeClass("error");
	$("[name='"+ elementName + "']").addClass("success");
	$("[name='"+ elementName + "-val-div']").text("");
	$("[name='"+ elementName + "-svg-error']").hide();
	$("[name='"+ elementName + "-svg-success']").show();
}

function addErrorToInput(elementName, errorString) {
	$("[name='"+ elementName + "']").removeClass("success");
	$("[name='"+ elementName + "']").addClass("error");
	$("[name='"+ elementName + "-val-div']").text(errorString);
	$("[name='"+ elementName + "-svg-success']").hide();
	$("[name='"+ elementName + "-svg-error']").show();
}

function validateRequiredNumeric(elementName) {
	var element = $("[name='"+ elementName + "']");
	if(element.val() >= 0 && element.val() < 2147483647) {
		addSuccessToInput(elementName);
		return true;
	}
	else {
		addErrorToInput(elementName, capitalizeFirstLetter(elementName) + " is required and cant be negative!")
		return false;
	}
}

function validateRequiredString(elementName) {
	var element = $("[name='"+ elementName + "']");	
	if(element.val().length > 0) {
		addSuccessToInput(elementName);
		return true;
	}
	else {
		addErrorToInput(elementName, "Please enter a value for " + capitalizeFirstLetter(elementName));
		return false;
	}
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}
