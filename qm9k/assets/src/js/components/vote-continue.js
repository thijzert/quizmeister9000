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

	if ( btn.disabled || btn.classList.contains("-disabled") ) {
		return;
	}

	let quizkey = mustSingle("main").dataset["quizkey"];
	let rv;

	if ( btn.classList.contains("-voted") ) {
		rv = await postJSON("vote-continue/"+quizkey, {vote: 0})
	} else {
		rv = await postJSON("vote-continue/"+quizkey, {vote: 1})
	}

	setVoteStatus(rv.MyVote);
}

export function setVoteStatus( myvote ) {
	myvote = !!myvote;
	all(".-js-vote-continue").forEach(e => {
		e.classList.toggle("-voted", myvote);
		e.classList.toggle("-check", myvote);
		e.classList.toggle("-alt", myvote);
	})
}

export function enableVoting( enabled ) {
	all(".-js-vote-continue").forEach(e => {
		e.classList.toggle("-disabled", !enabled);
		e.disabled = !enabled;
	})
}
