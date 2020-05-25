import { connect } from 'react-redux'
import Version from './Version.jsx'

const mapStateToProps = (state, ownProps) => ({
  version: ownProps.version
})

const VersionContainer = connect(mapStateToProps)(Version)

export default VersionContainer
