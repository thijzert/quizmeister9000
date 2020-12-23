
function post(url, obj) {
	let formData = objectToFormData( obj );
	formData.append("csrf_token", window._csrf);

	return fetch(url, {
		credentials: "same-origin",
		method: "post",
		body: formData
	});
}

export async function postJSON( url, obj ) {
	let response = null;
	try {
		response = await post( url, obj );
		return await response.json();
	} catch ( err ) {
		err.httpStatus = response?.status;
		err.url = url;
		err.formData = obj;
		throw err;
	}
}

export async function getJSON( url ) {
	let response = null;
	try {
		response = await fetch( url, { credentials: "same-origin" } );
		return await response.json();
	} catch ( err ) {
		err.httpStatus = response?.status;
		err.url = url;
		throw err;
	}
}

export async function getHTML( url ) {
	let response = null;
	try {
		response = await fetch(url, { credentials: 'same-origin' });

		if ( !response.ok ) {
			throw new Error("http request failed with status " + response.status + " " + response.statusText)
		}

		const text = await response.text();
		const dp = new DOMParser();
		const doc = dp.parseFromString(text, "text/html");

		return doc;
	} catch ( err ) {
		err.httpStatus = response?.status;
		err.url = url;
		throw err;
	}
}

export async function postHTML( url, obj ) {
	let response = null;
	try {
		response = await post( url, obj );

		if ( !response.ok ) {
			throw new Error("http request failed with status " + response.status + " " + response.statusText)
		}

		const text = await response.text();
		const dp = new DOMParser();
		const doc = dp.parseFromString(text, "text/html");

		return doc;
	} catch ( err ) {
		err.httpStatus = response?.status;
		err.url = url;
		err.formData = obj;
		throw err;
	}
}

// https://stackoverflow.com/questions/22783108/convert-js-object-to-form-data
// IE11 compatible
function objectToFormData(obj, rootName, ignoreList) {
	let formData = new FormData();

	function appendFormData(data, root) {
		if (!ignore(root)) {
			root = root || '';
			if (data instanceof File) {
				formData.append(root, data);
			} else if (Array.isArray(data)) {
				for (let i = 0; i < data.length; i++) {
					appendFormData(data[i], root + '[' + i + ']');
				}
			} else if (typeof data === 'object' && data) {
				for (let key in data) {
					// eslint-disable-next-line no-prototype-builtins
					if (data.hasOwnProperty(key)) {
						// eslint-disable-next-line max-depth
						if (root === '') {
							appendFormData(data[key], key);
						} else {
							appendFormData(data[key], root + `[${key}]`);
						}
					}
				}
			} else {
				if (data !== null && typeof data !== 'undefined') {
					formData.append(root, data);
				}
			}
		}
	}

	function ignore(root) {
		return Array.isArray(ignoreList) &&
			ignoreList.some(function (x) {
				return x === root;
			});
	}

	appendFormData(obj, rootName);

	return formData;
}
