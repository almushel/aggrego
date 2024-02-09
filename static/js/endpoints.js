async function createUser() {
	const nameInput = document.querySelector("#name-input");

	const requestBody = {
		name: nameInput.value
	};

	const response = await fetch("/v1/users", {
		headers: {
			"Content-Type": "application/json"
		},
		method: "POST",
		body: JSON.stringify(requestBody)
	});

	const rBody = await response.json();

	const success = await login(rBody.apikey);
	if (success) {
		localStorage.setItem(lsKeyName, rBody.apikey);
		nameInput.value = "";
	}
}

async function getUser(apikey) {
	const response = await fetch("/v1/users", {
		method: "GET",
		headers: {
			"Authorization": "ApiKey "+apikey
		},
	})
	if (!response.ok) {
		return false;
	}

	const result = await response.json();
	return result;
}

async function getPosts(apikey, offset, limit) {
	const response = await fetch(`/v1/posts?offset=${offset}&limit=${limit}`, {
		method: "GET",
		headers: {
			"Authorization": "ApiKey "+apikey
		}
	});
	const list = await response.json() || [];
	console.log(list);

	return {
		list,
		offset,
		totalCount: parseInt(response.headers.get("X-Total-Count")),
	}
}

async function getFeedFollows(apikey) {
	let response = await fetch("/v1/feed_follows", {
		method: "GET",
		headers: {
			"Authorization": "ApiKey "+apikey
		}
	});

	const result = await response.json() || [];
	return result;
}

async function getFeeds() {
	let response = await fetch("/v1/feeds", {
		method: "GET",
	});

	const result = await response.json();
	return result;
}

async function postFeed(name, url, apikey) {
	const body = {name, url};	
	const response = await fetch("/v1/feeds", {
		method: "POST",
		headers: {
			"Authorization": "ApiKey "+apikey
		},
		body: JSON.stringify(body)
	});

	return response.ok;
}

async function toggleFeedFollow(feed_id, follow) {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	const headers = {"Authorization": "ApiKey "+apikey};

	let response;
	if (follow) {
		response = await fetch("/v1/feed_follows", {
				method: "POST",
				headers,
				body: JSON.stringify({feed_id})
			}
		);
	} else {
		let ffID = "";
		const ffList = await getFeedFollows(apikey)
		for (const ff of ffList) {
			if (ff.feed_id == feed_id) {
				ffID = ff.id;
				break;
			}
		}

		if (!ffID) {
			console.log("toggleFeedFollow() can't delete feed_follow that doesn't exist");
			return false;
		} 

		response = await fetch("/v1/feed_follows/"+ffID, {
				method: "DELETE",
				headers
			}
		);
	}

	return true;
}

async function getLikes(apikey) {
	const response = await fetch ("/v1/post_likes", {
		method: "GET",
		headers: {
			"Authorization": "ApiKey "+apikey
		}
	});

	const result = await response.json() || [];
	return result;
}

async function likePost(apikey, post_id) {
	const body = {post_id};
	const response = await fetch("/v1/post_likes", {
		method: "POST",
		headers: {
			"Authorization": "ApiKey "+apikey
		},
		body: JSON.stringify(body)
	});

	const result = await response.json();
	return result;
}

async function unlikePost(apikey, likeID) {
	const response = await fetch("/v1/post_likes/"+likeID, {
		method: "DELETE",
		headers: {
			"Authorization": "ApiKey "+apikey
		}
	});

	return response.ok; 
}