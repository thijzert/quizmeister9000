
if ( document.readyState == "loading" ) {
	document.addEventListener( "DOMContentLoaded", main );
} else {
	main();
}

import { quizMain } from "./pages/quiz-viewer.js";
import { peerStatusMain } from "./components/peer-status.js";

function main() {
	peerStatusMain();

	let ndMain = document.querySelector("main")
	if ( ndMain ) {
		let c = ndMain.classList;
		if ( c.contains("quiz-viewer") ) {
			quizMain();
		}
	}
}
