import { closest, mustSingle, mustSingleRef, toggleIf } from "../lib/helpers.js";
import { getJSON } from "../lib/post.js";
import { enableVoting } from "../components/vote-continue.js";

export function quizMain() {
	let qcont = mustSingle(".-js-quiz-questions .question-container")

	setInterval( updateQuizStatus, 1200 );
	updateQuizStatus();
}

let currentRound = -2;

async function updateQuizStatus() {
	let quizkey = mustSingle("main").dataset["quizkey"];
	let status = await getJSON("/quiz-status/"+quizkey)

	enableVoting(status.VotingEnabled);

	if ( status.QuizStatus.Started && !status.QuizStatus.Grading && !status.QuizStatus.Finished ) {
		let qcont = mustSingle(".-js-quiz-questions .question-container")

		if ( status.CurrentRound.RoundNo != currentRound ) {
			// The round has changed. Update the form
			currentRound = status.CurrentRound.RoundNo;

			qcont.innerHTML = "";

			let template = mustSingle(".-js-template-answer");
			let tClass = "answer";
			if ( status.CurrentRound.ThisIsMe ) {
				template = mustSingle(".-js-template-question");
				tClass = "question";
			}
			status.CurrentRound.Questions.forEach(q => {
				let ndq = document.createElement("DIV");
				ndq.classList.add(tClass);
				ndq.innerHTML = template.innerHTML;

				if ( status.CurrentRound.ThisIsMe ) {
					mustSingleRef(ndq, "textarea").value = q.Question;
				} else {
					mustSingleRef(ndq, "-question").innerText = q.Question;
					mustSingleRef(ndq, "textarea").value = q.MyAnswer;
				}

				let txt = mustSingleRef(ndq, "textarea");
				txt.addEventListener("focus",questionFocus);
				txt.addEventListener("blur",questionBlur);

				qcont.appendChild(ndq);
				mustSingleRef(ndq, ".-number").innerText = qcont.children.length;
			});

			let title = mustSingle(".-js-quiz-questions h3")
			if ( status.CurrentRound.ThisIsMe ) {
				title.innerHTML = "This is you.";
			} else {
				title.innerHTML = "Please direct your attention to: ";
				let name = document.createElement("STRONG");
				name.innerText = status.CurrentRound.QuizMaster.Nick;
				title.appendChild(name);
			}
		}

		if ( status.CurrentRound.ThisIsMe ) {
			qcont.querySelectorAll(".question").forEach((ndq, i) => {
				const q = status.CurrentRound.Questions[i];
				if ( !q ) {
					return;
				}

				// TODO: check textareas for changes, and save
			})
		} else {
			qcont.querySelectorAll(".question").forEach((ndq, i) => {
				const q = status.CurrentRound.Questions[i];
				if ( !q ) {
					return;
				}

				mustSingleRef(ndq, "-question").innerText = q.Question;

				// TODO: check textareas for changes, and save
			})
		}
	}

	toggleIf( mustSingle(".-js-quiz-global-start"), !status.QuizStatus.Started )
	toggleIf( mustSingle(".-js-quiz-questions"), status.QuizStatus.Started && !status.QuizStatus.Grading && !status.QuizStatus.Finished )
	toggleIf( mustSingle(".-js-quiz-grading"), status.QuizStatus.Grading )
	toggleIf( mustSingle(".-js-quiz-global-end"), status.QuizStatus.Finished )


}

function questionFocus(e) {
	let nd = closest( e.target, ".question,.answer" );
	if ( nd ) {
		nd.classList.toggle("-focus", true);
	}
}
function questionBlur(e) {
	let nd = closest( e.target, ".question,.answer" );
	if ( nd ) {
		nd.classList.toggle("-focus", false);
	}
}
