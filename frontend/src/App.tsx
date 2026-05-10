import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Settings } from 'lucide-react';

type Node = {
  ID: number;
  IP: string;
  SSHPort: string;
  Status: string;
};

type Config = {
  flvx_account: string;
  flvx_password: string;
  flvx_api_url: string;
  cf_token: string;
  domain_name: string;
  check_api_url: string;
};

function App() {
  const [status, setStatus] = useState('Loading...');
  const [nodeStatus, setNodeStatus] = useState('Loading...');
  const [nodes, setNodes] = useState<Node[]>([]);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  
  const [formData, setFormData] = useState({ ip: '', ssh_port: '22', ssh_password: '' });
  const [configData, setConfigData] = useState<Config>({
    flvx_account: '', flvx_password: '', flvx_api_url: '', cf_token: '', domain_name: '', check_api_url: ''
  });

  const fetchStatus = () => {
    fetch('/api/status')
      .then(r => r.json())
      .then(d => {
        setStatus(d.status);
        setNodeStatus(d.node_status || 'Unknown');
      })
      .catch(() => {
        setStatus('Error');
        setNodeStatus('Error');
      });
  };

  const fetchNodes = () => {
    fetch('/api/nodes')
      .then(r => r.json())
      .then(d => setNodes(d))
      .catch(console.error);
  };

  const fetchConfig = () => {
    fetch('/api/config')
      .then(r => r.json())
      .then(d => setConfigData({
        flvx_account: d.flvx_account || '',
        flvx_password: d.flvx_password || '',
        flvx_api_url: d.flvx_api_url || '',
        cf_token: d.cf_token || '',
        domain_name: d.domain_name || '',
        check_api_url: d.check_api_url || '',
      }))
      .catch(console.error);
  };

  useEffect(() => {
    fetchStatus();
    fetchNodes();
    fetchConfig();
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

  const handleSaveConfig = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(configData),
      });
      if (res.ok) {
        setIsSettingsOpen(false);
        fetchConfig();
      } else {
        alert('Failed to save configuration');
      }
    } catch (err) {
      console.error(err);
      alert('Error saving configuration');
    }
  };

  return (
    <div className="container mx-auto p-8 max-w-4xl space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Flvx Monitor Dashboard</h1>
        <Dialog open={isSettingsOpen} onOpenChange={setIsSettingsOpen}>
          <DialogTrigger render={<Button variant="outline" size="icon" />}>
            <Settings className="h-4 w-4" />
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Settings & Configuration</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSaveConfig} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="flvx_api_url">Flvx API URL</Label>
                <Input 
                  id="flvx_api_url" 
                  value={configData.flvx_api_url} 
                  onChange={e => setConfigData({...configData, flvx_api_url: e.target.value})} 
                  placeholder="https://panel.flvx.com/api" 
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="flvx_account">Flvx Account</Label>
                <Input 
                  id="flvx_account" 
                  value={configData.flvx_account} 
                  onChange={e => setConfigData({...configData, flvx_account: e.target.value})} 
                  placeholder="admin@example.com" 
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="flvx_password">Flvx Password</Label>
                <Input 
                  id="flvx_password" 
                  type="password"
                  value={configData.flvx_password} 
                  onChange={e => setConfigData({...configData, flvx_password: e.target.value})} 
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="cf_token">Cloudflare API Token</Label>
                <Input 
                  id="cf_token" 
                  type="password"
                  value={configData.cf_token} 
                  onChange={e => setConfigData({...configData, cf_token: e.target.value})} 
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="domain_name">Domain Name (for DNS update)</Label>
                <Input 
                  id="domain_name" 
                  value={configData.domain_name} 
                  onChange={e => setConfigData({...configData, domain_name: e.target.value})} 
                  placeholder="node.example.com"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="check_api_url">3rd-Party GFW Check API URL</Label>
                <Input 
                  id="check_api_url" 
                  value={configData.check_api_url} 
                  onChange={e => setConfigData({...configData, check_api_url: e.target.value})} 
                  placeholder="https://api.itdog.cn/ping"
                />
              </div>
              <Button type="submit" className="w-full">Save Settings</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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
            <p className="text-sm text-gray-500 mt-2">Go Backend Monitor Daemon</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Flvx Node Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-semibold">
              {nodeStatus === 'Alive' ? (
                <span className="text-green-600">Alive</span>
              ) : nodeStatus === 'Dead' ? (
                <span className="text-red-600">Dead</span>
              ) : nodeStatus === 'Blocked' ? (
                <span className="text-orange-600">Blocked (GFW)</span>
              ) : nodeStatus === 'Replacing' ? (
                <span className="text-blue-600">Replacing...</span>
              ) : (
                <span className="text-gray-500">{nodeStatus}</span>
              )}
            </div>
            <p className="text-sm text-gray-500 mt-2">Active Node Connectivity</p>
          </CardContent>
        </Card>
      </div>

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
 App;
