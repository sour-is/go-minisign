package minisign

const VersionString="minisign 0.9"

const KeyNumBytes=8
const PasswordMaxBytes=1024
const SigAlg="Ed"
const SigAlgHashed="ED"
const KDFAlg="Sc"
const CHKAlg="B2"

const CommentPrefix="untrusted comment: "
const CommentMaxBytes=1024

const DefaultComment="signature from minisign secret key"
const SecretKeyDefaultComment="minisign encrypted secret key"

const TrustedCommentPrefix="trusted comment: "
const TrustedCommentMaxBytes=8192

const SigDefaultConfigDir=".minisign"
const SigDefaultConfigDirEnvVar="MINISIGN_CONFIG_DIR"
const SigDefaultPKFile="minisign.pub"
const SigDefaultSKFile="minisign.key"
const SigSuffix=".minisig"

