import { all, single, mustSingle } from "../lib/helpers.js";
import { getJSON } from "../lib/post.js";

export function peerStatusMain() {
	if ( single(".-js-update-peer-status") ) {
		setUpPeerStatus();
	}
}

function setUpPeerStatus() {
	window.setInterval( refreshPeerStatus, 222 );
}

async function refreshPeerStatus() {
	let quizkey = mustSingle("main").dataset["quizkey"];
	let peerStatus = await getJSON("/peer-status/"+quizkey)

	all(".-js-update-peer-status").forEach(elt => {
		peerStatus.Peers.forEach(peer => {
			let sel = ".-peer-id-" + peer.UserID;
			let peerElt = elt.querySelector(sel)
			if ( !peerElt ) {
				let tpl = elt.querySelector(".-js-template-peer");
				if ( !tpl ) {
					console.error("template not found");
					return;
				}
				peerElt = document.createElement("DIV")
				peerElt.classList.add("-peer")
				peerElt.classList.add(sel.substr(1));
				peerElt.innerHTML = tpl.innerHTML;
				elt.appendChild(peerElt);
			}

			let use = peerElt.querySelector("svg use");
			if ( use ) {
				use.setAttribute("xlink:href","#avatar-"+peer.Status);
				use.setAttribute("fill","#"+peer.Colour);
			}
			let nameElt = peerElt.querySelector(".-name");
			if ( nameElt ) {
				nameElt.textContent = peer.Nick;
			}
		})
	})
}

