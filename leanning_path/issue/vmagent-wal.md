I found vmagent is a high performance component of VictoriaMetrics. vmagent use in-memory queue to gather data block to flush data to remote address in most case.

Is it possible to lost data in vmagent if the data in memory queue and it's not flushed to remote address because I can't find the WAL component in vmagent?