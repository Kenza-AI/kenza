import { connect } from 'react-redux'

import JobDetails from './JobDetails.jsx'
import { fetchJobDetails } from './api'

const mapStateToProps = state => {
    const jobID = jobIDFromState(state)
    const projectID = projectIDFromState(state)
    const job = state.jobs.list.filter(job => job.id == jobID)[0]
    return {
        job: { ...job, id: jobID, project: { id: projectID } },
        fetchJobDetailsError: state.jobs.fetchJobDetailsError,
        isFetchingJobDetails: state.jobs.isFetchingJobDetails
    }
}

const mapDispatchToProps = (dispatch, ownProps) => ({
    fetchJobDetails: (projectID, jobID) => fetchJobDetails(dispatch, projectID, jobID, ownProps.config, ownProps.account)
})

const jobIDFromState = state => (
    state.location.pathname.split('/')[4]
)

const projectIDFromState = state => (
    state.location.pathname.split('/')[2]
)

const JobDetailsContainer = connect(mapStateToProps, mapDispatchToProps)(JobDetails)
export default JobDetailsContainer