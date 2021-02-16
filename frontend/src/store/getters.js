String.prototype.trim = function (char, type) {
  if (char) {
      if (type == 'left') {
          return this.replace(new RegExp('^\\'+char+'+', 'g'), '');
      } else if (type == 'right') {
          return this.replace(new RegExp('\\'+char+'+$', 'g'), '');
      }
      return this.replace(new RegExp('^\\'+char+'+|\\'+char+'+$', 'g'), '');
  }
  return this.replace(/^\s+|\s+$/g, '');
};

const getters = {
  isLogged: state => state.user !== null,
  isFiles: state => !state.loading && state.route.name === 'Files',
  isTrash: state => {
    let path = state.route.path.trim('/', 'right')
    return !state.loading && (path === '/files/.trash')
  },
  isListing: (state, getters) => getters.isFiles && state.req.isDir,
  isEditor: (state, getters) => getters.isFiles && (state.req.type === 'text' || state.req.type === 'textImmutable'),
  isPreview: state => state.previewMode,
  isSharing: state =>  !state.loading && state.route.name === 'Share',
  selectedCount: state => state.selected.length,
  progress : state => {
    if (state.upload.progress.length == 0) {
      return 0;
    }

    let sum = state.upload.progress.reduce((acc, val) => acc + val)
    return Math.ceil(sum / state.upload.size * 100);
  }
}

export default getters
