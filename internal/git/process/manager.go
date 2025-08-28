package process
m.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}


stdout, err := m.cmd.StdoutPipe()
if err != nil {
return err
}
stderr, err := m.cmd.StderrPipe()
if err != nil {
return err
}


if err := m.cmd.Start(); err != nil {
return err
}


// log stdout/stderr
go streamLog(stdout)
go streamLog(stderr)


return nil
}


func streamLog(r io.Reader) {
s := bufio.NewScanner(r)
for s.Scan() {
log.Println(s.Text())
}
}


func (m *Manager) GracefulShutdown(grace time.Duration) {
if m == nil || m.cmd == nil || m.cmd.Process == nil {
return
}
pid := m.cmd.Process.Pid
// send TERM to process group
syscall.Kill(-pid, syscall.SIGTERM)


c := make(chan error, 1)
go func() { c <- m.cmd.Wait() }()


select {
case <-c:
return
case <-time.After(grace):
// escalate
syscall.Kill(-pid, syscall.SIGKILL)
<-c
return
}
}