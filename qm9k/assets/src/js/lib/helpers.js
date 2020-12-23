
export function closest(elem, selector) {
	// Element.matches() polyfill
	if (!Element.prototype.matches) {
		Element.prototype.matches =
			Element.prototype.matchesSelector ||
			Element.prototype.mozMatchesSelector ||
			Element.prototype.msMatchesSelector ||
			Element.prototype.oMatchesSelector ||
			Element.prototype.webkitMatchesSelector ||
			function (s) {
				let matches = (this.document || this.ownerDocument).querySelectorAll(s);
				let i = matches.length;
				// eslint-disable-next-line no-empty
				while (--i >= 0 && matches.item(i) !== this) {}
				return i > -1;
			};
	}

	// Get the closest matching element
	for (; elem && elem !== document; elem = elem.parentNode) {
		if (elem.matches(selector)) return elem;
	}
	return null;
}

/**
 *
 * @param {node} parent
 * @param {event} evt
 * @param {childfilter} selector
 * @param {callback} handler
 */
export function on(parent, evt, selector, handler) {
	try {
		parent.addEventListener(evt, function (event) {
			// console.log( evt, event, event.target );
			try {
				let target = event.target;

				// IE11 - SVGElementInstance ellende
				if (target.correspondingUseElement) {
					target = target.correspondingUseElement;
				}

				if (target.matches(selector + ', ' + selector + ' *')) {
					handler.apply(target.closest(selector), arguments);
				}

			} catch (error) {
				// Nada
				console.log( 'on (1)', error, event.target, selector );
			}
		}, false);
	} catch (error) {
		console.log( 'on addEventListener', error );
	}
}

export function onClick( selector, handler ) {
	return on( document.body, "click", selector, handler );
}

export const all      = (c) => document.querySelectorAll(c);
export const single   = (c) => document.querySelector(c);

// Null-safe version of single()
export const mustSingle = (c) => document.querySelector(c) || document.createElement("TEMPLATE");

export const singleRef     = (n,c) => n.querySelector(c);
export const mustSingleRef = (n,c) => n.querySelector(c) || document.createElement("TEMPLATE");

// Toggle visibility of an element
export const toggleIf = ( elt, visible ) => {
	if ( !elt ) {
		return;
	}

	if ( visible ) {
		elt.style.display = null;
	} else {
		elt.style.display = "none";
	}
};

// Scroll de viewport zodat het element volledig zichtbaar is
export function scrollIntoViewIfNeeded( element ) {
	if ( !element ) {
		return;
	}

	let bounding = element.getBoundingClientRect();

	let viewportWidth  = window.innerWidth || document.documentElement.clientWidth;
	let viewportHeight = window.innerHeight || document.documentElement.clientHeight;

	if ( bounding.top >= 0 && bounding.left >= 0 && bounding.right <= viewportWidth && bounding.bottom <= viewportHeight ) {
		// Er hoeft niets gescrolld te worden
		return;
	}

	try {
		element.scrollIntoView({ behavior: "smooth", block: "nearest" });
	} catch ( _e ) {
		// Options niet ondersteund - scroll dan maar niet smooth.
		element.scrollIntoView( false );
	}
}

// Zet alle (ingevulde) velden uit dit element in een object.
// Inputs met namen als "foo.bar.baz" komen in sub-objecten, e.g. rv.foo.bar.baz = 1.
export function formToObject( form )
{
	let rv = {};
	for ( var x of form.querySelectorAll("input,select") )
	{
		let preamble = x.name.split(".");
		let k = preamble.pop();
		let v = null;

		if ( x.nodeName === "select" )
		{
			v = x.value;
		}
		else if ( x.type === "radio" )
		{
			v = x.checked ? x.value : null;
		}
		else if ( x.type === "checkbox" )
		{
			v = x.checked ? 1 : null;
		}
		else
		{
			v = x.value;
		}

		if ( !v )
		{
			// eslint-disable-next-line no-continue
			continue;
		}

		let cont = rv;
		for ( let kk of preamble )
		{
			if ( !cont[kk] )
			{
				cont[kk] = {};
			}
			cont = cont[kk];
		}
		cont[k] = v;
	}

	return rv;
}
