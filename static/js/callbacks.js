const lsKeyName = "aggrego-apikey";
const postPageSize = 25;

function onEnterCallback(e, callback) {
	if (e.key == "Enter") {
		callback();
	}
}

async function login(apikey) {
	const user = await getUser(apikey);
	if (!user) {
		return false;
	}

	document.getElementById("current-user").innerText = "Current user: " + user.name;
	document.getElementById("logged-in").style.display = "block";

	document.getElementById("logged-out").style.display = "none";

	loadPosts(apikey);
	updateFeedList(apikey);

	return true;
}

async function loginWithKeyInput() {
	const keyInput = document.getElementById("apikey-input");

	const success = await login(keyInput.value)
	if (success) {
		localStorage.setItem(lsKeyName, keyInput.value)
	}
}

function logout() {
	localStorage.removeItem(lsKeyName)

	document.getElementById("logged-in").style.display = "none";
	document.getElementById("posts-container").textContent = "";
	document.getElementById("logged-out").style.display = "block";
}

function loadPage(e) {
	e.preventDefault();
	const page = e.target.href;
	window.history.pushState({}, "", page);
	lsUpdatePosts();
}

function newPostElement(post) {
	const result = document.createElement("div")
	const title = document.createElement("h3");
	const url = document.createElement("a");
	url.href=post.url;
	url.innerText = post.title;

	const description = document.createElement("p");
	description.innerHTML = post.description;

	title.append(url);
	result.append(title, description);
	return result;
}

async function loadPosts(apikey, offset=0, limit=postPageSize) {
	const container = document.getElementById("posts-container");
	container.textContent = "";
	appendPosts(apikey, offset, limit);
}

async function appendPosts(apikey, offset=0, limit=postPageSize) {
	const container = document.getElementById("posts-container");
	
	const loadingElement = document.createElement("div");
	loadingElement.innerText = "Loading...";
	container.append(loadingElement);

	const posts = await getPosts(apikey, offset, limit);
	const children = [];
	if (posts.list.length > 0) {
		for (let post of posts.list) {
			children.push(newPostElement(post))
		}
		const nextOffset = posts.offset + children.length;

		const loadOnScrollEnd = async () => {
			const loadThreshold = 100;
			const rect = loadingElement.getBoundingClientRect();
			const viewDelta = rect.top - window.visualViewport.height;
			if (viewDelta <= loadThreshold) {
				document.removeEventListener("scrollend", loadOnScrollEnd)
				loadingElement.remove();
				await appendPosts(apikey, nextOffset);
			}
		};
		document.addEventListener("scrollend", loadOnScrollEnd)
	} else {
		loadingElement.textContent = `End of posts`;
	}

	container.append(...children, loadingElement)
}

async function updateFeedList(apikey) {
	const feeds = await getFeeds();
	const follows = await getFeedFollows(apikey);

	const followIDs = []
	for (let f of follows) {
		followIDs.push(f.feed_id)	
	}

	const feedList = document.getElementById("user-feeds");

	let count = 1;
	let elements = []
	for (let feed of feeds) {
		const url = new URL(feed.url);
		const elementID = `feed-list-item-${count}`;

		const e = document.createElement("li");
		e.innerHTML = 
			`<input type="checkbox" id="${elementID}" ${followIDs.includes(feed.id) ? "checked" : ""}>` +
			`<label for="${elementID}">${url.hostname}</label>`

		const checkbox = e.children[elementID];
		checkbox.addEventListener("input", async () => {
			const success = await toggleFeedFollow(feed.id, checkbox.checked)
			if (success) {
				lsUpdatePosts()
			} else {
				checkbox.checked = !checkbox.checked;
			}
		})

		elements.push(e)
		count++;
	}

	feedList.replaceChildren(...elements)
}

async function submitNewFeed() {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	const nameField = document.getElementById("new-feed-name");
	const urlField = document.getElementById("new-feed-url");

	if (nameField.value && urlField.value) {
		const success = await postFeed(nameField.value, urlField.value, apikey)
		if (success) {
			nameField.value = "";
			urlField.value = "";

			updateFeedList(apikey);
			loadPosts(apikey);
		}
	}
}

async function lsLogin() {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	login(apikey)
}

async function lsUpdatePosts() {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	loadPosts(apikey)
}