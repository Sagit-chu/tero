import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

type Node = {
  ID: number;
  IP: string;
  SSHPort: string;
  Status: string;
};

function App() {
  const [status, setStatus] = useState('Loading...');
  const [nodes, setNodes] = useState<Node[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [formData, setFormData] = useState({ ip: '', ssh_port: '22', ssh_password: '' });

  const fetchStatus = () => {
    fetch('/api/status')
      .then(r => r.json())
      .then(d => setStatus(d.status))
      .catch(() => setStatus('Error'));
  };

  const fetchNodes = () => {
    fetch('/api/nodes')
      .then(r => r.json())
      .then(d => setNodes(d))
      .catch(console.error);
  };

  useEffect(() => {
    fetchStatus();
    fetchNodes();
    const interval = setInterval(() => {
      fetchStatus();
      fetchNodes();
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleAddNode = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/nodes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });
      if (res.ok) {
        setIsDialogOpen(false);
        setFormData({ ip: '', ssh_port: '22', ssh_password: '' });
        fetchNodes();
      } else {
        alert('Failed to add node');
      }
    } catch (err) {
      console.error(err);
      alert('Error adding node');
    }
  };

  return (
    <div className="container mx-auto p-8 max-w-4xl space-y-8">
      <h1 className="text-3xl font-bold">Flvx Monitor Dashboard</h1>
      
      <Card>
        <CardHeader>
          <CardTitle>System Status</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-semibold">
            {status === 'ok' ? (
              <span className="text-green-600">Online</span>
            ) : status === 'Error' ? (
              <span className="text-red-600">Error</span>
            ) : (
              <span className="text-gray-500">Loading...</span>
            )}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Standby Nodes Pool</CardTitle>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger render={<Button />}>
              Add Node
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add Standby Node</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleAddNode} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="ip">IP Address</Label>
                  <Input 
                    id="ip" 
                    required 
                    value={formData.ip} 
                    onChange={e => setFormData({...formData, ip: e.target.value})} 
                    placeholder="1.1.1.1" 
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="port">SSH Port</Label>
                  <Input 
                    id="port" 
                    required 
                    value={formData.ssh_port} 
                    onChange={e => setFormData({...formData, ssh_port: e.target.value})} 
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="password">SSH Password</Label>
                  <Input 
                    id="password" 
                    type="password" 
                    required 
                    value={formData.ssh_password} 
                    onChange={e => setFormData({...formData, ssh_password: e.target.value})} 
                  />
                </div>
                <Button type="submit" className="w-full">Save Node</Button>
              </form>
            </DialogContent>
          </Dialog>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>IP Address</TableHead>
                <TableHead>SSH Port</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {nodes.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={3} className="text-center text-gray-500 py-8">
                    No nodes in the pool
                  </TableCell>
                </TableRow>
              ) : (
                nodes.map((node) => (
                  <TableRow key={node.ID}>
                    <TableCell className="font-medium">{node.IP}</TableCell>
                    <TableCell>{node.SSHPort}</TableCell>
                    <TableCell>
                      <span className={`px-2 py-1 rounded-full text-xs font-semibold ${
                        node.Status === 'standby' ? 'bg-blue-100 text-blue-800' :
                        node.Status === 'active' ? 'bg-green-100 text-green-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {node.Status}
                      </span>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}

export default App;
