import { fetchURL, fetchJSON, removePrefix } from './utils'

export async function list() {
    return fetchJSON('/api/favorites')
}

export async function remove(hash) {
    const res = await fetchURL(`/api/favorite/${hash}`, {
        method: 'DELETE'
    })

    if (res.status !== 200) {
        throw new Error(res.status)
    }
}

export async function create(url) {
    url = removePrefix(url) + `/api/favorite${url}`
    return fetchJSON(url, {
        method: 'POST',
    })
}