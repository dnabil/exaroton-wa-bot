/**
 * Checks if the given input array has an even length.
 * @param {array} elements - The array to check (even: the input elements, odd: the helper (<small>) elements.)
 * @returns {boolean} - True if the length is even, false if not
 */
function validateInput(elements) {
  if (elements.length % 2 != 0) {
    throw new Error(
      "validateInput: The elements array must contain pairs of input and helper elements."
    );
  }

  let valid = true;
  for (let i = 0; i < elements.length; i += 2) {
    const input = elements[i];
    const helper = elements[i + 1];

    valid = _validateInputHelper(input, helper);
  }

  return valid;
}

/**
 * Attaches input event listeners to each input element and its corresponding helper element for validation
 * with a delay to prevent immediate validation feedback.
 * @param {Array} elements - Array containing input elements and their corresponding helper elements.
 */
function setupValidationListeners(elements) {
  if (elements.length % 2 !== 0) {
    throw new Error(
      "setupValidationListeners: The elements array must contain pairs of input and helper elements."
    );
  }

  // prevent immediate feedback for web autocomplete.
  setTimeout(() => {
    for (let i = 0; i < elements.length; i += 2) {
      const inputElement = elements[i];
      const helperElement = elements[i + 1];
      let timeoutId;

      inputElement.addEventListener("input", () => {
        clearTimeout(timeoutId);

        timeoutId = setTimeout(() => {
          console.debug("validating...");
          _validateInputHelper(inputElement, helperElement);
        }, 500);
      });
    }
  }, 500);
}

/**
 * Validates a single input element with its corresponding helper element.
 * @param {HTMLInputElement} inputElement - The input element to validate.
 * @param {HTMLSmallElement} helperElement - The corresponding helper element to display the validation message.
 * @returns {boolean} - True if the input is valid, false if not.
 * @private
 */
function _validateInputHelper(inputElement, helperElement) {
  helperElement.textContent = "";

  const isValid = inputElement.checkValidity();
  if (!isValid) {
    console.debug("not valid");
    console.debug(" validation message:", inputElement.validationMessage);
    helperElement.textContent = inputElement.validationMessage;
  }

  inputElement.setAttribute("aria-invalid", !isValid);

  return isValid;
}
