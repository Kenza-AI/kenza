import React, { Fragment } from 'react'
import { closeIcon } from '../icons'

class AddProject extends React.Component {

  componentDidUpdate(prevProps) {
    const pendingProject = this.props.pendingProject

    if (pendingProject.accessToken !== prevProps.pendingProject.accessToken) {
      fetchReposIfNeeded(pendingProject, this.props.fetchReposHandler)
      return
    } else {
      if (pendingProject.repo === undefined && pendingProject.repos && pendingProject.repos.length > 0) {
        this.props.pendingProjectChangeHandler(pendingProject.title, pendingProject.description, pendingProject.repos[0].full_name, '')
      }
    }
  }

  render() {
    const {
      pendingProject,
      createProjectHandler,
      hideAddProjectHandler,
      pendingProjectChangeHandler,
    } = this.props

    return (
      <section className='add-project'>
        <div className='fullscreen-overlay' onClick={hideAddProjectHandler}></div>
        <div className='fullscreen-overlay-content'>
          <button className='dismiss-modal' onClick={hideAddProjectHandler}>{closeIcon()}</button>
          <h5>Create New Project</h5>
          {accessToken(pendingProject, pendingProjectChangeHandler)}
          {branch(pendingProject, pendingProjectChangeHandler)}
          {repo(pendingProject, pendingProjectChangeHandler)}
          {title(pendingProject, pendingProjectChangeHandler)}
          {description(pendingProject, pendingProjectChangeHandler)}
          {createProjectCTA(pendingProject, createProjectHandler)}
        </div>
      </section>
    )
  }
}

const accessToken = (pendingProject, pendingProjectChangeHandler) => (
  <Fragment>
    <div className='form-wrapper'>
      <input
        onChange={e => pendingProjectChangeHandler({ ...pendingProject, accessToken: e.target.value })}
        type='text'
        defaultValue={pendingProject && pendingProject.accessToken || ''}
        placeholder='Access token with read permissions on the repo' required />
      <label>* Required</label>
    </div>
    <span className='form-error'></span>
  </Fragment>
)

const title = (pendingProject, pendingProjectChangeHandler) => (
  <Fragment>
    <div className='form-wrapper'>
      <input
        onChange={e => pendingProjectChangeHandler({ ...pendingProject, title: e.target.value })}
        type='text'
        defaultValue={pendingProject && pendingProject.title || ''}
        placeholder='Give a name to this project' required />
      <label>* Required</label>
    </div>
    <span className='form-error'></span>
  </Fragment>
)

const description = (pendingProject, pendingProjectChangeHandler) => {
  return (
    <Fragment>
      <div className='form-wrapper'>
        <input
          onChange={e => pendingProjectChangeHandler({ ...pendingProject, description: e.target.value })}
          type='text'
          defaultValue={pendingProject && pendingProject.description || ''}
          placeholder='Give a description of this project...' />
      </div>
      <span className='form-error'></span>
    </Fragment>
  )
}

const repo = (pendingProject, pendingProjectChangeHandler) => (
  <Fragment>
    <div className='form-wrapper'>
      <input
        onChange={e => pendingProjectChangeHandler({ ...pendingProject, repo: e.target.value })}
        type='text'
        defaultValue={pendingProject && pendingProject.repo || ''}
        placeholder='Repository clone URL' required />
      <label>* Required</label>
    </div>
    <span className='form-error'></span>
  </Fragment>
)

const branch = (pendingProject, pendingProjectChangeHandler) => (
  <Fragment>
    <div className='form-wrapper'>
      <input
        onChange={e => pendingProjectChangeHandler({ ...pendingProject, branch: e.target.value })}
        type='text'
        defaultValue={pendingProject && pendingProject.branch || ''}
        placeholder='Git ref regex e.g. (refs/heads/.*) to build all branches' required />
      <label>* Required</label>
    </div>
    <span className='form-error'></span>
  </Fragment>
)

const createProjectCTA = (pendingProject, createProjectHandler) => {
  return (
    <div id='add-project-footer'>
      <button onClick={() => validateAndCreate(
        pendingProject,
        createProjectHandler)}>Create
        </button>
    </div>
  )
}

const validateAndCreate = (pendingProject, createProjectHandler) => {

  let isTitleValid = true
  let isBranchValid = true
  let isRepoValid = true

  // if (title.trim() === '') {
  //   isTitleValid = false
  //   titleError.textContent = 'Cannot be left blank'
  // }

  // if (!isTitleValid) {
  //   titleInputWrapper.classList.add('error')
  // } else {
  //   titleError.textContent = ''
  //   titleInputWrapper.classList.remove('error')
  // }

  if (isTitleValid && isRepoValid && isBranchValid) {
    createProjectHandler(pendingProject)
  }
}

const fetchReposIfNeeded = (pendingProject, fetchReposHandler) => {
  if (pendingProject && pendingProject.gitHub_access_token) {
    fetchReposHandler(pendingProject.gitHub_access_token)
  }
}

export default AddProject
