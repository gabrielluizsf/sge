//go:build gots_overlay

// Overlay for connectionStateLocked: Go's, but for the verified chains of
// a client, worked out once asked for (see denoChains). Should they not
// verify in Go, the state holds none.

package tls

func (c *Conn) connectionStateLocked() ConnectionState {
	if c.isClient {
		c.denoChains()
	}
	var state ConnectionState
	state.HandshakeComplete = c.isHandshakeComplete.Load()
	state.Version = c.vers
	state.NegotiatedProtocol = c.clientProtocol
	state.DidResume = c.didResume
	state.HelloRetryRequest = c.didHRR
	state.testingOnlyPeerSignatureAlgorithm = c.peerSigAlg
	state.CurveID = c.curveID
	state.NegotiatedProtocolIsMutual = true
	state.ServerName = c.serverName
	state.CipherSuite = c.cipherSuite
	state.PeerCertificates = c.peerCertificates
	state.LocalCertificate = c.localCertificate
	state.VerifiedChains = c.verifiedChains
	state.SignedCertificateTimestamps = c.scts
	state.OCSPResponse = c.ocspResponse
	if (!c.didResume || c.extMasterSecret) && c.vers != VersionTLS13 {
		if c.clientFinishedIsFirst {
			state.TLSUnique = c.clientFinished[:]
		} else {
			state.TLSUnique = c.serverFinished[:]
		}
	}
	if c.config.Renegotiation != RenegotiateNever {
		state.ekm = noEKMBecauseRenegotiation
	} else if c.vers != VersionTLS13 && !c.extMasterSecret {
		state.ekm = noEKMBecauseNoEMS
	} else {
		state.ekm = c.ekm
	}
	state.ECHAccepted = c.echAccepted
	return state
}
