import { closest, on, mustSingle, all } from "../lib/helpers.js";
import { postJSON } from "../lib/post.js";

export function voteContinueMain() {
	on( document, "click", ".-js-vote-continue", voteContinue)
}

async function voteContinue(e) {
	const btn = closest( e.target, ".-js-vote-continue" );
	if ( !btn ) {
		return;
	}

	let quizkey = mustSingle("main").dataset["quizkey"];
	let rv;

	if ( btn.classList.contains("-voted") ) {
		rv = await postJSON("/vote-continue/"+quizkey, {vote: 0})
	} else {
		rv = await postJSON("/vote-continue/"+quizkey, {vote: 1})
	}

	let tog = !!rv.MyVote;
	all(".-js-vote-continue").forEach(e => {
		e.classList.toggle("-voted", tog);
		e.classList.toggle("-check", tog);
	})
}
