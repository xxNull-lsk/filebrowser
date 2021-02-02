import { fetchURL, fetchJSON } from './utils'

export async function list() {
    return fetchJSON('/api/favorites')
}

export async function get(path) {
    return fetchJSON(`/api/favorite${path}`)
}

export async function remove(hash) {
    const res = await fetchURL(`/api/favorite/${hash}`, {
        method: 'DELETE'
    })

    if (res.status !== 200) {
        throw new Error(res.status)
    }
}

export async function create(path, name) {
    console.error("function create")
    return fetchJSON('/api/favorite', {
        method: 'POST',
        body: JSON.stringify({
            path: path,
            name: name
        }
        )
    })
}